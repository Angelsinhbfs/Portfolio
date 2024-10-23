import React from "react";
import '../styles/Tile.css'

interface PRT_TileProps {
    label: string;
    onClick?: () => void;
    callback?: (details:string) => void;
    imgURL?: string;
    tileDetails: string;
}
const PTile: React.FC<PRT_TileProps> = ({label, onClick, callback, tileDetails}) => {
    const handleClick = () => {
        if (onClick) {
            onClick();
        }
        if (callback) {
            callback(tileDetails);
        }
    };
    return (
        <div className="tile-bg" onClick={handleClick}>
            <div className="tile-overlay">
                <span className="overlay-text">{label}</span>
            </div>
        </div>
    );
};

export default PTile;