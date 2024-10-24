import React, {FormEvent, useEffect, useState} from 'react';
import Modal from "./Modal";
import MarkdownPreview from "./MarkdownPreview";
import "../styles/Modal.css"
import fetchWithToken from "../../fetchWithToken";

interface UploadModalProps {
    isOpen: boolean;
    onRequestClose: () => void;
    origin: string;
    update: PortfolioEntry | null;
}

const UploadModal: React.FC<UploadModalProps> = ({ isOpen, onRequestClose, origin, update }) => {
    const [title, setTitle] = useState<string>('');
    const [tags, setTags] = useState<string>('');
    const [body, setBody] = useState<string>('');
    const [id, setId] = useState<string>('');

    useEffect(() => {
        if (update) {
            setTitle(update.label || ''); // Initialize title with update.title if available
            setTags(update.tags ? update.tags.join(' ') : ''); // Initialize tags with update.tags if available
            setBody(update.details || ''); // Initialize body with update.body if available
            setId(update._id || '');
        }
    }, [update]);

    const uploadImage = async (blob: Blob): Promise<string> => {
        const formData = new FormData();
        formData.append('image', blob);

        try {
            const response = await fetchWithToken(`${origin}portfolio/img/`, {
                method: 'POST',
                body: formData,
            });
            return response.url;
        } catch (error) {
            console.error('Error uploading image:', error);
            return '';
        }
    };

    const handleImagePaste = async (e: React.ClipboardEvent<HTMLTextAreaElement>) => {
        const items = e.clipboardData?.items;
        if (items) {
            let newBody = body; // Initialize newBody with the current body content
            for (let i = 0; i < items.length; i++) {
                if (items[i].type.indexOf('image') !== -1) {
                    e.preventDefault(); // Prevent the default paste behavior

                    const blob = items[i].getAsFile();
                    if (blob) {
                        const startPos = e.currentTarget?.selectionStart || 0;
                        const endPos = e.currentTarget?.selectionEnd || 0;
                        const imageUrl = await uploadImage(blob);


                        const imageTag = `![Pasted Image](${imageUrl})`;
                        newBody = newBody.slice(0, startPos) + imageTag + newBody.slice(endPos);
                    }
                }
            }
            setBody(newBody); // Update the body content with all pasted images
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        try {
            const response = await fetchWithToken(`${origin}portfolio/create/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    'id': id,
                    'label': title,
                    'tags': tags.split(' '),
                    'details': body,
                }),
            });
            setTitle('');
            setTags('');
            setBody('');
            setId('');
        } catch (error) {
            console.error('Error uploading tile data:', error);
        }

        onRequestClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={onRequestClose} isAuth={true}>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Title:</label>
                    <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
                    <label>Tags:</label>
                    <input type="text" value={tags} onChange={(e) => setTags(e.target.value)} required />
                </div>
                <div className="modal-content">
                    <div className="modal-input">
                        <label>Input</label>
                        <textarea  value={body} onChange={(e) => setBody(e.target.value)} onPaste={handleImagePaste} required />
                    </div>
                    <MarkdownPreview markdown={body} /> {/* Render the MarkdownPreview component */}
                </div>
                <button type="submit">Submit</button>
            </form>
        </Modal>
    );
};

export default UploadModal;