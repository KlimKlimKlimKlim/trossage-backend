package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// registerUser registers a new user
//
//	@Summary	Register user
//	@Tags		auth
//	@Param		request	body		dto.RegisterUserRequest	true	"Registration data"
//	@Success	201		{object}	dto.SuccessUserAndTokenResponse
//	@Failure	400		{object}	dto.ErrorResponse	"Invalid request data"
//	@Failure	409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure	500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router		/auth/register [post]
func (h *handler) registerUser(ctx *gin.Context) {
	var req dto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	user, jwtPair, err := h.service.CreateUser(ctx, req.Login, req.Password, req.DisplayName)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	middleware.SetUserID(ctx, user.ID)

	var resp dto.UserAndTokenResponse
	resp.Fill(user, jwtPair.AccessToken, jwtPair.RefreshToken)

	response.Created(ctx, resp)
}
