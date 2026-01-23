package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// deleteCurrentUser deletes current user account
//
//	@Summary		Delete current user
//	@Description	Permanently deletes authenticated user account and all associated data. Requires password confirmation.
//	@Tags			users
//	@Security		BearerAuth
//	@Param			request	body		dto.DeleteUserRequest	true	"Password confirmation"
//	@Success		200		{object}	dto.SuccessEmptyResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request body"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized or invalid password"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/users/me [delete]
func (h *handler) deleteCurrentUser(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	var req dto.DeleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	if err := h.service.DeleteUser(ctx, userID, req.Password); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.OK(ctx, nil)
}
