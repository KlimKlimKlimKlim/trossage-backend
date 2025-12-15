package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// logoutUser revokes current refresh token (single session)
// @Summary      Logout
// @Description  Revokes current refresh token and ends session on this device
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} dto.LogoutResponse
// @Failure      401 {object} dto.ErrorResponse "Invalid or expired refresh token"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /auth/logout [post]
func (s *state) logoutUser(ctx *gin.Context) {
	tokenID, ok := middlewares.GetTokenID(ctx)
	if !ok {
		response.HandleError(ctx, errors.New("token_id is empty"))
		return
	}

	if err := s.controller.Logout(ctx, tokenID); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.OK(ctx, nil)
}
