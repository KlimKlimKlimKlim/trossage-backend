package api

import (
	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/dto"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

func (s *state) loginUser(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	accessToken, refreshToken, err := s.controller.LoginUser(ctx, req.Login, req.Password)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.TokenResponse
	resp.Fill(accessToken, refreshToken)

	response.Created(ctx, resp)
}
