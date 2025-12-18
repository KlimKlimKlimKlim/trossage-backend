package dto

import (
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

// RegisterUserRequest contains registration data
type RegisterUserRequest struct {
	Login       string `json:"login"        binding:"required,min=3,max=20" example:"john_doe"`
	Password    string `json:"password"     binding:"required,min=8,max=63" example:"securePass123"`
	DisplayName string `json:"display_name" binding:"required,max=20"       example:"John Doe"`
}

// LoginUserRequest contains user credentials
type LoginUserRequest struct {
	Login    string `json:"login"    binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"securePass123"`
}

// UserResponse contains user information
type UserResponse struct {
	ID          int64     `json:"id"           example:"12345"`
	Login       string    `json:"login"        example:"john_doe"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	CreatedAt   time.Time `json:"created_at"   example:"2025-12-15T09:00:00Z"`
}

func (dto *UserResponse) Fill(user models.User) {
	dto.ID = user.ID
	dto.Login = user.Login
	dto.DisplayName = user.DisplayName
	dto.CreatedAt = user.CreatedAt
}

// TokenResponse contains access token pair
type TokenResponse struct {
	AccessToken  string `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func (dto *TokenResponse) Fill(accessToken, refreshToken string) {
	dto.AccessToken = accessToken
	dto.RefreshToken = refreshToken
}

// UserAndTokenResponse contains user data with tokens
type UserAndTokenResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

func (dto *UserAndTokenResponse) Fill(user models.User, accessToken, refreshToken string) {
	dto.User.Fill(user)
	dto.Token.Fill(accessToken, refreshToken)
}

// UpdateUserRequest contains fields that can be updated
type UpdateUserRequest struct {
	DisplayName string `json:"display_name,omitempty" binding:"omitempty,min=1,max=20" example:"John Doe Updated"`
	OldPassword string `json:"old_password,omitempty" binding:"omitempty"              example:"currentPassword123"`
	NewPassword string `json:"new_password,omitempty" binding:"omitempty,min=8,max=63" example:"newSecurePass123"`
}

// UpdateUserResponse represents update response data
type UpdateUserResponse struct {
	User                UserResponse `json:"user"`
	TokensRevoked       bool         `json:"tokens_revoked"                  example:"true"`
	TokensRevokedReason string       `json:"tokens_revoked_reason,omitempty" example:"password changed"`
}

func (dto *UpdateUserResponse) Fill(user models.User, tokensRevoked bool, tokensRevokedReason string) {
	dto.User.Fill(user)
	dto.TokensRevoked = tokensRevoked
	dto.TokensRevokedReason = tokensRevokedReason
}

// DeleteUserRequest contains password confirmation for account deletion
type DeleteUserRequest struct {
	Password string `json:"password" binding:"required" example:"currentPassword123"`
}

type UsersSearchResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"  example:"42"`
	Limit  int            `json:"limit"  example:"20"`
	Offset int            `json:"offset" example:"0"`
}

func (dto *UsersSearchResponse) Fill(users []models.User, total, limit, offset int) {
	dto.Users = make([]UserResponse, len(users))
	for i, user := range users {
		dto.Users[i].Fill(user)
	}
	dto.Total = total
	dto.Limit = limit
	dto.Offset = offset
}
