package users

import "time"

// User represents the database entity
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never export this
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterRequest defines the payload for registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the payload for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is what we send back on successful login
type AuthResponse struct {
	Token string `json:"token"`
}
