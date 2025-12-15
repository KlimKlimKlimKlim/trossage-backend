package dto

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

// RegisterUserRequest contains registration data
type RegisterUserRequest struct {
	Login       string `json:"login" binding:"required,min=3,max=20" example:"john_doe"`
	Password    string `json:"password" binding:"required,min=8,max=63" example:"securePass123"`
	DisplayName string `json:"display_name" binding:"required,max=20" example:"John Doe"`
}

// LoginUserRequest contains user credentials
type LoginUserRequest struct {
	Login    string `json:"login" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"securePass123"`
}

// UserResponse contains user information
type UserResponse struct {
	ID          int64  `json:"id" example:"12345"`
	Login       string `json:"login" example:"john_doe"`
	DisplayName string `json:"display_name" example:"John Doe"`
}

func (dto *UserResponse) Fill(user models.User) {
	dto.ID = user.ID
	dto.Login = user.Login
	dto.DisplayName = user.DisplayName
}

// TokenResponse contains access token pair
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func (dto *TokenResponse) Fill(accessToken, refreshToken string) {
	dto.AccessToken = accessToken
	dto.RefreshToken = refreshToken
}

// CreateUserResponse contains created user data with tokens
type CreateUserResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

func (dto *CreateUserResponse) Fill(user models.User, accessToken, refreshToken string) {
	dto.User.Fill(user)
	dto.Token.Fill(accessToken, refreshToken)
}
