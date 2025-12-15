package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// refreshToken issues new token pair using refresh token
// @Summary      Refresh tokens
// @Description  Generates new access and refresh tokens using valid refresh token
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} dto.RefreshTokenResponse
// @Failure      401 {object} dto.ErrorResponse "Invalid or expired refresh token"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /auth/refresh [post]
func (s *state) refreshToken(ctx *gin.Context) {
	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, errors.New("user_id is empty"))
		return
	}

	tokenID, ok := middlewares.GetTokenID(ctx)
	if !ok {
		response.HandleError(ctx, errors.New("token_id is empty"))
		return
	}

	accessToken, refreshTokenNew, err := s.controller.RefreshToken(ctx, userID, tokenID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.TokenResponse
	resp.Fill(accessToken, refreshTokenNew)

	response.OK(ctx, resp)
}
