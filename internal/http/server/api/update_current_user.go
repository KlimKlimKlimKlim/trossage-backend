package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/dto"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/middleware"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/response"
)

// updateCurrentUser updates current user profile
//
//	@Summary		Update current user
//	@Description	Updates authenticated user display name and/or password. Password change revokes all refresh tokens.
//	@Tags			users
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateUserRequest	true	"Fields to update"
//	@Success		200		{object}	dto.SuccessUpdateUserResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request data or same password"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/users/me [patch]
func (h *handler) updateCurrentUser(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	if err := req.Validate(); err != nil {
		response.HandleError(ctx, err)
		return
	}

	user, tokenRevocation, err := h.service.UpdateUser(
		ctx,
		userID,
		req.DisplayName,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.UpdateUserResponse
	resp.Fill(user, tokenRevocation.Revoked, tokenRevocation.Reason)

	response.OK(ctx, resp)
}
