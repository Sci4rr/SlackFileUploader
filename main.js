import axios from 'axios';

const fileInput = document.querySelector('#file-input');
const uploadButton = document.querySelector('#upload-button');
const statusElement = document.querySelector('#status');

function updateStatus(message, isError = false) {
    statusElement.textContent = message;
    statusElement.style.color = isError ? 'red' : 'black';
    logToConsole(`Status Update: ${message}`, isError ? 'error' : 'info');
}

function logToConsole(message, level = 'info') {
    const timestamp = new Date().toISOString();
    if (level === 'error') {
        console.error(`[${timestamp}] Error: ${message}`);
    } else {
        console.log(`[${timestamp}] ${message}`);
    }
}

async function uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);

    logToConsole('Attempting to upload file...');
    try {
        const response = await axios.post(`${process.env.BACKEND_URL}/upload`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
        });

        updateStatus(`Upload successful: ${response.data.message}`);
        logToBoth('File uploaded successfully.', 'info');
    } catch (error) {
        console.error('Error uploading file:', error);
        updateStatus('Error uploading file. Please try again.', true);
    }
}

uploadButton.addEventListener('click', (e) => {
    e.preventDefault();
    
    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];
        updateStatus('Uploading file...');
        uploadFile(file);
    } else {
        updateStatus('Please select a file to upload.', true);
    }
});