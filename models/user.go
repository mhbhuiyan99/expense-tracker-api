package models

// User represents a user in the system
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"` // never expose password in JSON responses
	CreatedAt string `json:"created_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" valid:"Required"`
	Email    string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required;MinSize(6)"`
}

// LoginRequest is the expected JSON body for POST /api/v1/auth/login
type LoginRequest struct {
	Email    string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required"`
}

// LoginData is the data returned on successful login
type LoginData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
