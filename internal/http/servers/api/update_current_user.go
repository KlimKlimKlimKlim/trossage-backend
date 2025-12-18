package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
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
func (s *state) updateCurrentUser(ctx *gin.Context) {
	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	if req.DisplayName == "" && req.NewPassword == "" {
		response.HandleError(ctx, derrors.ErrEmptyBody)
		return
	}

	if req.NewPassword != "" && req.OldPassword == "" {
		response.HandleError(ctx, derrors.ErrUnauthorized)
		return
	}

	user, tokensRevoked, tokensRevokedReason, err := s.controller.UpdateUser(
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
	resp.Fill(user, tokensRevoked, tokensRevokedReason)

	response.OK(ctx, resp)
}
