package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// logoutUserAll revokes all user's refresh tokens (all sessions)
// @Summary      Logout from all devices
// @Description  Revokes all refresh tokens and ends all user sessions
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} dto.LogoutAllResponse
// @Failure      401 {object} dto.ErrorResponse "Unauthorized"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /auth/logout-all [post]
func (s *state) logoutUserAll(ctx *gin.Context) {
	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, errors.New("user_id is empty"))
		return
	}

	if err := s.controller.LogoutAll(ctx, userID); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.OK(ctx, nil)
}
