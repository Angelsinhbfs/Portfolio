import React, { useState, ChangeEvent, FormEvent } from 'react';
import Modal from "./Modal";
import MarkdownPreview from "./MarkdownPreview";
import "../styles/Modal.css"

interface UploadModalProps {
    isOpen: boolean;
    onRequestClose: () => void;
    origin: string;
}

const UploadModal: React.FC<UploadModalProps> = ({ isOpen, onRequestClose, origin }) => {
    const [title, setTitle] = useState<string>('');
    const [tags, setTags] = useState<string>('');
    const [body, setBody] = useState<string>('');
    const [image, setImage] = useState<File | null>(null);
    const [imageUrl, setImageUrl] = useState<string>('');

    const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setImage(e.target.files[0]);
        }
    };

    const handleImageUpload = async () => {
        if (!image) return;

        const formData = new FormData();
        formData.append('image', image);

        try {
            const response = await fetch(`${origin}portfolio/img`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: formData,
            });
            if (!response.ok){
                throw new Error('Network response was not ok')
            }
            const data = await response.json();

            setImageUrl(data.url);
        } catch (error) {
            console.error('Error uploading image:', error);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        try {
            const response = await fetch(`${origin}portfolio/create`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    'label': title,
                    'tags': tags.split(' '),
                    'details': body,
                }),
            });
            if (!response.ok){
                throw new Error('Network response was not ok')
            }
            const data = await response.json();

            setImageUrl(data.url);
        } catch (error) {
            console.error('Error uploading tile data:', error);
        }

        onRequestClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={onRequestClose}>
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
                        <textarea  value={body} onChange={(e) => setBody(e.target.value)} required />
                    </div>
                    <MarkdownPreview markdown={body} /> {/* Render the MarkdownPreview component */}
                </div>
                <div>
                    <label>Image:</label>
                    <input type="file" accept="image/*" onChange={handleImageChange} />
                    <button type="button" onClick={handleImageUpload}>Upload Image</button>
                </div>
                {imageUrl && <img src={imageUrl} alt="Uploaded" style={{ maxWidth: '100%' }} />}
                <button type="submit">Submit</button>
            </form>
        </Modal>
    );
};

export default UploadModal;