import axios from 'axios';
require('dotenv').config();
const BASE_URL = process.env.BACKEND_URL;
async function uploadFile(fileData) {
    try {
        const formData = new FormData();
        formData.append('file', fileData);
        const response = await axios.post(`${BASE_URL}/upload-file`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        console.log(response.data);
        return response.data; 
    } catch (error) {
        console.error("Error uploading file:", error);
        throw error;
    }
}
async function getFileStatus(fileId) {
    try {
        const response = await axios.get(`${BASE_URL}/file-status/${fileId}`);
        console.log(response.data);
        return response.data; 
    } catch (error) {
        console.error("Error getting file status:", error);
        throw error;
    }
}
export { uploadFile, getFileStatus };