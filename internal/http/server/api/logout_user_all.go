package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// logoutUserAll revokes all user's refresh tokens (all sessions)
//
//	@Summary		Logout from all devices
//	@Description	Revokes all refresh tokens and ends all user sessions
//	@Tags			auth
//	@Security		BearerAuth
//	@Success		200	{object}	dto.SuccessEmptyResponse
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/logout-all [post]
func (h *handler) logoutUserAll(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	if err := h.service.LogoutAll(ctx, userID); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.OK(ctx, nil)
}
