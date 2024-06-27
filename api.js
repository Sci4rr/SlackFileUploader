import axios from 'axios';
require('dotenv').config();

const api = axios.create({
    baseURL: process.env.BACKEND_URL,
    headers: { 'Content-Type': 'multipart/form-data' }
});

function logMessage(type, message) {
    const prefix = type.toUpperCase();
    console.log(`${prefix}: ${message}`);
}

async function uploadFile(fileData) {
    try {
        const formData = new FormData();
        formData.append('mfile', fileData);
        const response = await api.post('/upload-file', formData);
        logMessage('info', `File uploaded successfully: ${JSON.stringify(response.data)}`);
        return response.data;
    } catch (error) {
     logMessage('error', `Error uploading file: ${error}`);
        throw error;
    }
}

async function getFileStatus(fileId) {
    try {
        const response = await api.get(`/file-status/${fileId}`);
        logMessage('info', `File status retrieved: ${JSON.stringify(response.data)}`);
        return response.data;
    } catch (error) {
        logMessage('error', `Error getting file status: ${error}`);
        throw error;
    }
}

export { uploadFile, getFileStatus };