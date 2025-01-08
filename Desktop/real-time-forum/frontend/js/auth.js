// Handle Registration
// Handle Registration
const registrationForm = document.getElementById('registration-form');
if (registrationForm) {
    registrationForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(registrationForm);
        const data = Object.fromEntries(formData);

        try {
            const response = await fetch('/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            });

            if (response.ok) {
                alert('Registration successful!');
                window.location.href = 'index.html'; // Redirect to the home page
            } else {
                try {
                    const error = await response.json();
                    alert(`Registration failed: ${error.message}`);
                } catch (err) {
                    alert('An unexpected error occurred during registration.');
                    console.error('Failed to parse server error:', err);
                }
            }
        } catch (err) {
            alert('An error occurred during registration. Please try again.');
            console.error(err);
        }
    });
}


// Handle Login
const loginForm = document.getElementById('login-form');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = new FormData(loginForm);
        const data = Object.fromEntries(formData);

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            });

            //if (response.ok) {
                //alert('Login successful!');
                //window.location.href = 'posts.html'; // Redirect to posts page after login
            //} else {
                //const error = await response.json();
                //alert(`Login failed: ${error.message}`);
            //}
            if (response.ok) {
                alert('Login successful!');
                localStorage.setItem('user', JSON.stringify(await response.json())); // Save user info
                window.location.href = 'homepage.html'; // Redirect to the homepage
            }
        } catch (err) {
            alert('An error occurred during login. Please try again.');
            console.error(err);
        }
    });
}
