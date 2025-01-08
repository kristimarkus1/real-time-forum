package main

import (
	"time"
)

// Post represents a post in the forum
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllPosts retrieves all posts from the database
