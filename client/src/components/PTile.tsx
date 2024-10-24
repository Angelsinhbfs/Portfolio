import React from "react";
import '../styles/Tile.css'

interface PRT_TileProps {
    label: string;
    onClick?: () => void;
    callback?: (details:string, tileId:string) => void;
    imgURL?: string;
    tileDetails: string;
    tileId: string;
}
const PTile: React.FC<PRT_TileProps> = ({label, onClick, callback, tileDetails, tileId}) => {
    const handleClick = () => {
        if (onClick) {
            onClick();
        }
        if (callback) {
            callback(tileDetails, tileId);
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