package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// logoutUser revokes current refresh token (single session)
//
//	@Summary		Logout
//	@Description	Revokes current refresh token and ends session on this device
//	@Tags			auth
//	@Security		BearerAuth
//	@Success		200	{object}	dto.SuccessEmptyResponse
//	@Failure		401	{object}	dto.ErrorResponse	"Invalid or expired refresh token"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/auth/logout [post]
func (s *state) logoutUser(ctx *gin.Context) {
	tokenID, ok := middlewares.GetTokenID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrTokenIDIsEmpty)
		return
	}

	if err := s.controller.Logout(ctx, tokenID); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.OK(ctx, nil)
}
