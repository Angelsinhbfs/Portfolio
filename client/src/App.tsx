import React, {useEffect, useState} from 'react';
import Modal from './components/Modal';
import PTile from './components/PTile';
import UploadModal from "./components/UploadModal";
import "./styles/Main.css"
import LoginButton from "./components/LoginButton";
import Md2HTML from "./Md2HTML";
import "./localTypes"

const App: React.FC = () => {
    const [isModalOpen, setModalOpen] = useState(false);
    const [isUploadOpen, setUploadOpen] = useState(false);
    const [tileInfo, setTileInfo] = useState('');
    const [tileList,setTileList] = useState<PortfolioEntry[]>([]);
    const [origin,assignOrigin] = useState('')

    const openModal = () => setModalOpen(true);
    const closeModal = () => setModalOpen(false);
    const openUpload = () => setUploadOpen(true);
    const closeUpload = () => setUploadOpen(false);
    const getDetails = (details:string) => setTileInfo(details)
    const setOrigin = (origin:string) => assignOrigin(origin)

    const fetchTileList = async () => {
        try {
            const response = await fetch(`${origin}portfolio/load`);
            if (response.ok) {
                const data = await response.json();
                setTileList(data); // Assuming the response data is an array of tiles
                return data; // Return the data
            } else {
                throw new Error('Failed to fetch data');
            }
        } catch (error) {
            console.error('Error fetching tile list:', error);
            throw error; // Re-throw the error
        }
    };
    useEffect(() => {
        const currentUrl = window.location.href;
        const url = new URL(currentUrl);
        const baseUrl = `${url.origin}${url.pathname}`;
        setOrigin(baseUrl)
        fetchTileList().then();
    }, []); // Empty dependency array to run the effect only once when the component mounts

    return (
        <div>
            <h1>Sidney Fernandez</h1> <h2>Software Engineer | Game Dev | XR Enthusiast | Maker</h2>
            <div className="top-right">
                <LoginButton></LoginButton>
                <button onClick={openUpload}>Upload</button>
            </div>
            <p>About Me: <br/>
                I am a passionate Software Engineer specializing in Unity3D and XR development. With a career that includes working at NASA's Jet Propulsion Laboratory, I have had the privilege of contributing to projects (InSight, Curiosity, Perseverance) that push the boundaries of technology and exploration. My work has been featured in National Geographic, and I am honored to have received the NASA Software of the Year award. <br/><br/>

                Beyond my professional achievements, I am an avid enthusiast of tabletop role-playing games (TTRPGs), live-action role-playing (LARP), building and flying FPV drones, RC robotics, and game design. These hobbies not only fuel my creativity but also keep me engaged with the latest technological advancements. <br/><br/>

                As a lifelong learner, I value curiosity and continuous learning above all else. I am always eager to explore new ideas and technologies, and I strive to bring innovation and excellence to every project I undertake.<br/><br/>

                Below you will find more information about some of my work. I hope you enjoy ^_^
            </p>
            <div className="tile-field" id="tile field">
                {tileList.map((tile, index) => (
                    <PTile key={index} label={tile.label} onClick={openModal} callback={getDetails} tileDetails={tile.details} />
                ))}
            </div>

            <UploadModal isOpen={isUploadOpen} onRequestClose={closeUpload} origin={`${origin}`}/>
            <Modal isOpen={isModalOpen} onClose={closeModal}>
                <div dangerouslySetInnerHTML={{__html: Md2HTML(tileInfo)}}/>
            </Modal>
        </div>
    );
};

export default App;