import React, { useState } from 'react';
import axios from 'axios';

const FileUploadForm = () => {
    const [selectedFiles, setSelectedFiles] = useState([]);

    const handleFileChange = (event) => {
        setSelectedFiles(event.target.files);
    };

    const handleFileUpload = async () => {
        if (!selectedFiles.length) {
            alert('Please select one or more files first.');
            return;
        }

        const uploadUrl = `${process.env.REACT_APP_API_URL}/YOUR_API_ENDPOINT`;

        const uploadPromises = Array.from(selectedFiles).map(file => {
            const formData = new FormData();
            formData.append('file', file);
            
            return axios.post(uploadUrl, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
        });

        try {
            await Promise.all(uploadPromises);
            alert('All files uploaded successfully');
        } catch (error) {
            console.error('Error uploading files', error);
            alert('Error uploading some or all files');
        }
    };

    return (
        <div>
            <input type="file" onChange={handleFileChange} multiple />
            <button onClick={handleFileUpload}>Upload</button>
        </div>
    );
};

export default FileUploadForm;