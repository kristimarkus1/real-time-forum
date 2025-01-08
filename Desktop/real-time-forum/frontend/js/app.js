document.addEventListener('DOMContentLoaded', () => {
    const registerButton = document.getElementById('register-btn');
    const loginButton = document.getElementById('login-btn');

    // Redirect to registration page
    registerButton.addEventListener('click', () => {
        window.location.href = 'registration.html';
    });

    // Redirect to login page
    loginButton.addEventListener('click', () => {
        window.location.href = 'login.html';
    });
});
