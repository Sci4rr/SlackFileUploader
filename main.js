import axios from 'axios';
const fileInput = document.querySelector('#file-input');
const uploadButton = document.querySelector('#upload-button');
const statusElement = document.querySelector('#query');
function updateStatus(message, isError = false) {
  statusElement.textContent = message;
  if (isError) {
    statusElement.style.color = 'red';
  } else {
    statusElement.style.color = 'black';
  }
}
async function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  try {
    const response = await axios.post(process.env.BACKEND_URL + '/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    updateStatus(`Upload successful: ${response.data.message}`);
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