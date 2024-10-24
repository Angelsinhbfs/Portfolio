import React, {useEffect, useState} from 'react';
import Modal from './components/Modal';
import PTile from './components/PTile';
import UploadModal from "./components/UploadModal";
import "./styles/Main.css"
import Md2HTML from "./Md2HTML";
import "./localTypes"
import fetchWithToken from "../fetchWithToken";

const App: React.FC = () => {
    const [isModalOpen, setModalOpen] = useState(false);
    const [isUploadOpen, setUploadOpen] = useState(false);
    const [tileInfo, setTileInfo] = useState('');
    const [tileId, setTileId] = useState('');
    const [tileList,setTileList] = useState<PortfolioEntry[]>([]);
    const [origin,assignOrigin] = useState('')
    const [isLoggedIn, setIsLoggedIn] = useState(false); // State to track user login status
    const [toUpdate, setToUpdate] = useState<PortfolioEntry | null>(null)

    const openModal = () => setModalOpen(true);
    const closeModal = () => setModalOpen(false);
    const openUpload = () => setUploadOpen(true);
    const edit = (toEdit:string) => {
        const filteredEntry = tileList.find(entry=>entry._id === toEdit);
        if (filteredEntry) {
            setToUpdate(filteredEntry);
            return;
        }
        setToUpdate(null);
    }
    const closeUpload = () => {
        setUploadOpen(false);
        fetchTileList().then();
    }
    const getDetails = (details:string, tileId:string) => {
        setTileInfo(details);
        setTileId(tileId);
    }
    const setOrigin = (origin:string) => assignOrigin(origin)

    const deleteEntry = async () => {
        try {
            await fetchWithToken(`${origin}portfolio/delete/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    'id': tileId,
                }),
            });
            await fetchTileList();
            closeModal();
        } catch (error) {
            console.error('Error fetching tile list:', error);
            throw error; // Re-throw the error
        }
    }

    const editEntry = () => {
        edit(tileId);
        closeModal();
        openUpload();
    }

    const fetchTileList = async () => {
        try {
            const response = await fetch(`${origin}portfolio/load/`);
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

    // Function to check if the user is logged in based on the presence of the token cookie
    const checkLoginStatus = () => {
        // Check if the token cookie is set
        const tokenCookie = sessionStorage.getItem('token');
        setIsLoggedIn(tokenCookie != null);
    };

    // Function to handle login button click
    const handleLogin = async () => {
        try {
            const response = await fetch(`${origin}login/`, {
                method: 'POST',
            });

            if (!response.ok) {
                throw new Error('Could not authenticate');
            }

            // Extract the token from the response headers
            const token = response.headers.get('token');

            if (token) {
                // Set the token as a cookie
                sessionStorage.setItem('token', token);
                //setIsLoggedIn(true);
                checkLoginStatus();
            } else {
                throw new Error('Token not found in response headers');
            }
        } catch (error) {
            console.error('Error handling login:', error);
            // Handle the error accordingly
        }
    };
    useEffect(() => {
        const currentUrl = window.location.href;
        const url = new URL(currentUrl);
        const baseUrl = `${url.origin}${url.pathname}`;
        setOrigin(baseUrl)
        fetchTileList().then();
        checkLoginStatus();
    }, []); // Empty dependency array to run the effect only once when the component mounts

    return (
        <div>
            <h1>Sidney Fernandez</h1>
            <h2>Software Engineer | Game Dev | XR Enthusiast | Maker</h2>
            <div className="top-right">
                {isLoggedIn ? ( // Conditionally render based on user login status
                    <button onClick={openUpload}>Update</button>
                ) : (
                    <button onClick={handleLogin}>Login</button>
                )}
            </div>
            <p>About Me: <br/>
                I am a passionate Software Engineer specializing in Unity3D and XR development. With a career that includes working at NASA's Jet Propulsion Laboratory, I have had the privilege of contributing to projects (InSight, Curiosity, Perseverance) that push the boundaries of technology and exploration. My work has been featured in National Geographic, and I am honored to have received the NASA Software of the Year award. <br/><br/>

                Beyond my professional achievements, I am an avid enthusiast of tabletop role-playing games (TTRPGs), live-action role-playing (LARP), building and flying FPV drones, RC robotics, and game design. These hobbies not only fuel my creativity but also keep me engaged with the latest technological advancements. <br/><br/>

                As a lifelong learner, I value curiosity and continuous learning above all else. I am always eager to explore new ideas and technologies, and I strive to bring innovation and excellence to every project I undertake.<br/><br/>

                Below you will find more information about some of my work. I hope you enjoy ^_^
            </p>
            <div className="tile-field" id="tile field">
                {tileList && tileList.map((tile, index) => (
                    <PTile key={index} label={tile.label} onClick={openModal} callback={getDetails} tileDetails={tile.details} tileId={tile._id} />
                ))}
            </div>

            <UploadModal isOpen={isUploadOpen} onRequestClose={closeUpload} origin={`${origin}`} update={toUpdate}/>
            <Modal isOpen={isModalOpen} onClose={closeModal} onDelete={deleteEntry} isAuth={isLoggedIn} onEdit={editEntry}>
                <div dangerouslySetInnerHTML={{__html: Md2HTML(tileInfo)}}/>
            </Modal>
        </div>
    );
};

export default App;