package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// getCurrentUser returns current authenticated user data
//
//	@Summary		Get current user
//	@Description	Returns authenticated user profile information
//	@Tags			users
//	@Security		BearerAuth
//	@Success		200	{object}	dto.SuccessUserResponse
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/users/me [get]
func (h *handler) getCurrentUser(ctx *gin.Context) {
	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.UserResponse
	resp.Fill(user)

	response.OK(ctx, resp)
}
