package main

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        int    `json:"id" db:"id"`
	Nickname  string `json:"nickname" db:"nickname"`
	Age       string `json:"age" db:"age"`
	Gender    string `json:"gender" db:"gender"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
}

// CreateUser adds a new user to the database
func CreateUser(user *User) error {
	// Hash the password
	log.Println("Hashing password for user:", user.Nickname)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the user into the database
	log.Println("Inserting user into database")
	_, err = db.Exec(
		"INSERT INTO users (nickname, age, gender, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, hashedPassword,
	)
	if err != nil {
		log.Println("Error inserting user into database:", err)
	}
	return err
}

// AuthenticateUser checks if the provided credentials are valid
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	err := GetDB().QueryRow("SELECT id, nickname, email, password FROM users WHERE nickname = ? OR email = ?", username, username).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
