import React from 'react';

const Footer = () => {
    return (
        <footer>
            <p>&copy; {new Date().getFullYear()} SlackFileUploader</p>
            <ul>
                <li><a href="#about">About Us</a></li>
                <li><a href="#help">Help</a></li>
                <li><a href="#privacy">Privacy Policy</a></li>
                <li><a href="#terms">Terms of Service</a></li>
            </ul>
        </footer>
    );  
};

export default Footer;