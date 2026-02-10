package dto

import (
	"time"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/model"
)

type RegisterUserRequest struct {
	Login       string `json:"login"        binding:"required,min=3,max=20" example:"john_doe"`
	Password    string `json:"password"     binding:"required,min=8,max=63" example:"securePass123"`
	DisplayName string `json:"display_name" binding:"required,max=20"       example:"John Doe"`
}

type LoginUserRequest struct {
	Login    string `json:"login"    binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"securePass123"`
}

type UserResponse struct {
	CreatedAt   time.Time `json:"created_at"   example:"2025-12-15T09:00:00Z"`
	Login       string    `json:"login"        example:"john_doe"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	ID          int64     `json:"id"           example:"12345"`
}

func (dto *UserResponse) Fill(user model.User) {
	dto.ID = user.ID
	dto.Login = user.Login
	dto.DisplayName = user.DisplayName
	dto.CreatedAt = user.CreatedAt
}

type UserAndTokenResponse struct {
	Token TokenResponse `json:"token"`
	User  UserResponse  `json:"user"`
}

func (dto *UserAndTokenResponse) Fill(user model.User, accessToken, refreshToken string) {
	dto.User.Fill(user)
	dto.Token.Fill(accessToken, refreshToken)
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name,omitempty" binding:"omitempty,min=1,max=20" example:"John Doe Updated"`
	OldPassword string `json:"old_password,omitempty" binding:"omitempty"              example:"currentPassword123"`
	NewPassword string `json:"new_password,omitempty" binding:"omitempty,min=8,max=63" example:"newSecurePass123"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.DisplayName == "" && r.NewPassword == "" {
		return derrors.ErrEmptyBody
	}

	if r.NewPassword != "" && r.OldPassword == "" {
		return derrors.ErrUnauthorized
	}

	return nil
}

type UpdateUserResponse struct {
	TokensRevokedReason string       `json:"tokens_revoked_reason,omitempty" example:"password changed"`
	User                UserResponse `json:"user"`
	TokensRevoked       bool         `json:"tokens_revoked"                  example:"true"`
}

func (dto *UpdateUserResponse) Fill(user model.User, tokensRevoked bool, tokensRevokedReason string) {
	dto.User.Fill(user)
	dto.TokensRevoked = tokensRevoked
	dto.TokensRevokedReason = tokensRevokedReason
}

type DeleteUserRequest struct {
	Password string `json:"password" binding:"required" example:"currentPassword123"`
}

type UsersSearchResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"  example:"42"`
	Limit  int            `json:"limit"  example:"20"`
	Offset int            `json:"offset" example:"0"`
}

func (dto *UsersSearchResponse) Fill(users []model.User, total, limit, offset int) {
	dto.Users = make([]UserResponse, len(users))
	for i, user := range users {
		dto.Users[i].Fill(user)
	}

	dto.Total = total
	dto.Limit = limit
	dto.Offset = offset
}
