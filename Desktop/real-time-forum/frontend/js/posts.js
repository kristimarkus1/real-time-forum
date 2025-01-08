document.addEventListener('DOMContentLoaded', async () => {
    const postsList = document.getElementById('posts-list');
    const postForm = document.getElementById('post-form');
    const categorySelect = document.getElementById('category-select');

    // Fetch the current user's data from localStorage
    const userData = JSON.parse(localStorage.getItem('user'));

    // Directly update the username element without creating unnecessary variables
    if (userData) {
        const usernameElement = document.getElementById('username');
        if (usernameElement) {
            usernameElement.textContent = userData.nickname || userData.email;
        } else {
            console.error("Username element not found in the DOM.");
        }
    } else {
        alert('You are not logged in!');
        window.location.href = 'login.html';
        return;
    }

    // Fetch and display posts
    async function fetchPosts(category = 'All') {
        try {
            const response = await fetch(`/posts?category=${category}`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            });
    
            if (response.ok) {
                const posts = await response.json();
                postsList.innerHTML = ''; // Clear previous posts
                console.log('Posts fetched:', posts);
    
                posts.forEach(post => {
                    const li = document.createElement('li');
                    li.innerHTML = `
                        <div class="post">
                            <h3>${post.title}</h3>
                            <p>${post.content}</p>
                            <p><em>Category: ${post.category}</em></p>
                            <button data-id="${post.id}" class="delete-button">Delete</button>
                        </div>
                    `;
                    postsList.appendChild(li);
                });
            } else {
                const error = await response.json();
                console.error('Failed to fetch posts:', error);
                alert(`Failed to fetch posts: ${error.message}`);
            }
        } catch (err) {
            console.error('Error fetching posts:', err);
            alert('An error occurred while fetching posts.');
        }
    }

    // Event listener to handle category changes
    categorySelect.addEventListener('change', (e) => {
        const selectedCategory = e.target.value;
        fetchPosts(selectedCategory);
    });

    // Fetch posts on page load
    await fetchPosts();

    // Handle post creation
    postForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(postForm);
        const postData = {
            title: formData.get('title'),
            content: formData.get('content'),
            category: formData.get('category') || 'All',  // Default to "All" if no category is selected
        };

        console.log("Posting the following data:", postData); // Debugging line to verify payload

        try {
            const response = await fetch('/posts', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(postData),
            });

            if (response.ok) {
                await fetchPosts();
                postForm.reset();
                alert('Post created successfully!');
            } else {
                const error = await response.json();
                console.error('Failed to create post:', error);
                alert(`Failed to create post: ${error.message}`);
            }
        } catch (err) {
            console.error('Error creating post:', err);
            alert('An error occurred while creating the post.');
        }
    });

    // Handle post deletion
    postsList.addEventListener('click', async (e) => {
        if (e.target.classList.contains('delete-button')) {
            const postID = e.target.getAttribute('data-id');
            if (!confirm('Are you sure you want to delete this post?')) return;

            try {
                const response = await fetch(`/posts?id=${postID}`, {
                    method: 'DELETE',
                });

                if (response.ok) {
                    await fetchPosts();
                    alert('Post deleted successfully!');
                } else {
                    const error = await response.json();
                    console.error('Failed to delete post:', error);
                    alert(`Failed to delete post: ${error.message}`);
                }
            } catch (err) {
                console.error('Error deleting post:', err);
                alert('An error occurred while deleting the post.');
            }
        }
    });

    await fetchPosts();
});
