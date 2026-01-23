package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// refreshToken issues new token pair using refresh token
//
//	@Summary		Refresh tokens
//	@Description	Generates new access and refresh tokens using valid refresh token
//	@Tags			auth
//	@Security		BearerAuth
//	@Success		200	{object}	dto.SuccessTokenResponse
//	@Failure		401	{object}	dto.ErrorResponse	"Invalid or expired refresh token"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *handler) refreshToken(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	tokenID, ok := middleware.GetTokenID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrTokenIDIsEmpty)
		return
	}

	accessToken, refreshTokenNew, err := h.service.RefreshToken(ctx, userID, tokenID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.TokenResponse
	resp.Fill(accessToken, refreshTokenNew)

	response.OK(ctx, resp)
}
