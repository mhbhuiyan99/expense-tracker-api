package services

import (
	"expense-tracker-api/models"
	"fmt"
)

// RegisterUser validates the registration request and creates a new user
// Returns an error if the email already exists
func RegisterUser(req models.RegisterRequest) error {
	existing, err := models.GetUserByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}	
	if existing != nil {
		return fmt.Errorf("email already exists")
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
	}
	return models.CreateUser(user)
}

// LoginUser validates credentials and returns login data
// Returns nil and an error if credentials are invalid
func LoginUser(req models.LoginRequest) (*models.LoginData, error) {
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing users: %w", err)
	}
	if user == nil || user.Password != req.Password {
		return nil, fmt.Errorf("invalid email or password")
	}

	return &models.LoginData{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}