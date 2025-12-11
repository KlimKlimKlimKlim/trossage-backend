package dto

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type RegisterUserRequest struct {
	Login       string `json:"login" binding:"required,min=3,max=20"`
	Password    string `json:"password" binding:"required,min=8,max=63"`
	DisplayName string `json:"display_name" binding:"required,max=20"`
}

type LoginUserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

func (dto *UserResponse) Fill(user models.User) {
	dto.ID = user.ID
	dto.Login = user.Login
	dto.DisplayName = user.DisplayName
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (dto *TokenResponse) Fill(accessToken, refreshToken string) {
	dto.AccessToken = accessToken
	dto.RefreshToken = refreshToken
}

type CreateUserResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

func (dto *CreateUserResponse) Fill(user models.User, accessToken, refreshToken string) {
	dto.User.Fill(user)
	dto.Token.Fill(accessToken, refreshToken)
}
