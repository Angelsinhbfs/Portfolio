import React, {useEffect} from "react";
import ReactDOM from "react-dom";
import '../styles/Modal.css'

interface ModalProps {
    isOpen: boolean;
    isAuth: boolean;
    onClose: () => void;
    onDelete?: () => void;
    onEdit?: () => void;
    children: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, isAuth, onClose, onDelete, onEdit, children }) => {
    useEffect(()=>{
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                onClose();
            }
        }
        if (isOpen) {
            document.addEventListener('keydown', handleKeyDown);
        }

        // Cleanup the event listener when the component is unmounted or when isOpen changes
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [isOpen, onClose])
    if (!isOpen) return null;

    return ReactDOM.createPortal(
        <div className="modal-overlay">
            <div className="modal-content">
                <button onClick={onClose} className="close-button">X</button>
                {children}
                {isAuth ? (
                    <div>
                    <button onClick={onDelete}>Delete</button>
                    <button onClick={onEdit}>Edit</button>
                    </div>
                ):(<button onClick={onClose} className="close-button">X</button>)}
            </div>
        </div>,
        document.body
    );
};


export default Modal;