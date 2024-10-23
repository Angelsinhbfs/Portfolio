import React, {useEffect} from "react";
import ReactDOM from "react-dom";
import '../styles/Modal.css'

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    children: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, children }) => {
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
            </div>
        </div>,
        document.body
    );
};


export default Modal;