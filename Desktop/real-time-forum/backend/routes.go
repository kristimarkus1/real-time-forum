package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterHandler called")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Reading JSON body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received payload: %s\n", string(body))

	var user struct {
		Nickname  string `json:"nickname"`
		Age       string `json:"age"`
		Gender    string `json:"gender"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	log.Println("Decoding JSON body")
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	log.Printf("Parsed user data: %+v\n", user)

	// Validate input fields
	log.Println("Validating input fields")
	if user.Nickname == "" || user.Email == "" || user.Password == "" {
		log.Println("Validation failed: Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	if user.Age == "" {
		log.Println("Validation failed: Age is required")
		log.Printf("Received Age: %v", user.Age)
		http.Error(w, "Age is required", http.StatusBadRequest)
		return
	}

	log.Println("Creating user in database")
	err = CreateUser(&User{
		Nickname:  user.Nickname,
		Age:       user.Age,
		Gender:    user.Gender,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		log.Println("Error creating user in database:", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Println("User registration successful")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	user, err := AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		log.Println("Authentication failed:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	log.Printf("User authenticated: %+v\n", user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"nickname": user.Nickname,
		"email":    user.Email,
	})
}

// PostsHandler handles creating, retrieving, and deleting posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "All" // Default to "All" if no category is provided
	}

	switch r.Method {
	case http.MethodGet:
		log.Println("Fetching posts by category:", category)
		posts, err := GetAllPosts(category)
		if err != nil {
			log.Println("Error fetching posts:", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(posts)

	case http.MethodPost:
		log.Println("Creating a new post")
		var post Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid data format", http.StatusBadRequest)
			return
		}

		if post.Title == "" || post.Content == "" || post.Category == "" {
			log.Println("Validation failed: Missing required fields")
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		err = CreatePost(&post)
		if err != nil {
			log.Println("Error creating post:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		log.Println("Post created successfully")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)

	case http.MethodDelete:
		log.Println("Processing delete request")
		postID := r.URL.Query().Get("id")
		if postID == "" {
			log.Println("Missing post ID")
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}

		err := DeletePost(postID)
		if err != nil {
			log.Println("Error deleting post:", err)
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		log.Println("Post deleted successfully, ID:", postID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})

	default:
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// GetAllPosts retrieves all posts from the database, optionally filtering by category
func GetAllPosts(category string) ([]Post, error) {
	query := "SELECT id, title, content, category, created_at FROM posts"
	if category != "All" {
		query += " WHERE category = ?"
	}

	rows, err := GetDB().Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	log.Println("Fetched posts:", posts)
	return posts, nil
}

// CreatePost inserts a new post into the database
func CreatePost(post *Post) error {
	result, err := GetDB().Exec("INSERT INTO posts (title, content, category) VALUES (?, ?, ?)",
		post.Title, post.Content, post.Category)
	if err != nil {
		log.Println("Error inserting post:", err)
		return err
	}

	lastID, _ := result.LastInsertId()
	log.Println("Post created successfully with ID:", lastID)
	return nil
}

// DeletePost removes a post from the database by its ID
func DeletePost(postID string) error {
	query := "DELETE FROM posts WHERE id = ?"
	result, err := GetDB().Exec(query, postID)
	if err != nil {
		log.Println("Error executing delete query:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("No post found with the given ID:", postID)
		return errors.New("no post found with the given ID")
	}

	log.Println("Post deleted successfully, ID:", postID)
	return nil
}
