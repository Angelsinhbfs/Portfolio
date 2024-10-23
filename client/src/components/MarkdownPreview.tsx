import React from "react";
import Md2HTML from "../Md2HTML";

const MarkdownPreview: React.FC<{ markdown: string }> = ({ markdown }) => {
    return (
        <div className="modal-preview">
            <div dangerouslySetInnerHTML={{__html: Md2HTML(markdown)}}/>
        </div>
    );
};

export default MarkdownPreview;
