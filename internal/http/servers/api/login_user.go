package api

import (
	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/dto"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// loginUser authenticates existing user
// @Summary      Login user
// @Tags         auth
// @Param        request body dto.LoginUserRequest true "Login credentials"
// @Success      200 {object} dto.LoginUserResponse
// @Failure      400 {object} dto.ErrorResponse "Invalid request data"
// @Failure      401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure      500 {object} dto.ErrorResponse "Internal server error"
// @Router       /auth/login [post]
func (s *state) loginUser(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	user, accessToken, refreshToken, err := s.controller.LoginUser(ctx, req.Login, req.Password)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.UserAndTokenResponse
	resp.Fill(user, accessToken, refreshToken)

	response.OK(ctx, resp)
}
