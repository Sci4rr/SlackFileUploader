import axios from 'axios';
require('dotenv').config();
const BASE_URL = process.env.BACKEND_URL;

function logMessage(type, message) {
    switch (type) {
        case 'info':
            console.log(`INFO: ${message}`);
            break;
        case 'error':
            console.error(`ERROR: ${message}`);
            break;
        default:
            console.log(message);
    }
}

async function uploadFile(fileData) {
    try {
        const formData = new FormData();
        formData.append('file', fileData);
        const response = await axios.post(`${BASE_URL}/upload-file`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        logMessage('info', `File uploaded successfully: ${JSON.stringify(response.data)}`);
        return response.data;
    } catch (error) {
                logMessage('error', `Error uploading file: ${error}`);
        throw error;
    }
}

async function getFileStatus(fileId) {
    try {
        const response = await axios.get(`${BASE_URL}/file-status/${fileId}`);
        logMessage('info', `File status retrieved: ${JSON.stringify(response.data)}`);
        return response.data;
    } catch (error) {
        logMessage('error', `Error getting file status: ${error}`);
        throw error;
    }
}

export { uploadFile, getFileStatus };