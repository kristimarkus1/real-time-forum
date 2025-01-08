package main

import (
	"database/sql"
	"log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitializeDB initializes the database connection and creates schema if needed
func InitializeDB() {
	// Check if the database file exists
	if _, err := os.Stat("forum.db"); os.IsNotExist(err) {
		log.Println("Database not found, creating a new database...")

		// Create the database file and apply the schema
		cmd := exec.Command("sqlite3", "forum.db", ".read migrations/schema.sql")
		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to execute schema.sql:", err)
		}
		log.Println("Database schema applied successfully.")
	}

	// Connect to the database
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully.")
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		return err
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Database connected successfully.")
	return nil
}
