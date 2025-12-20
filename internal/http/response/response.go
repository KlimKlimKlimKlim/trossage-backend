package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
)

type Response struct {
	Data      any    `json:"data"`
	Error     string `json:"error"`
	IsSuccess bool   `json:"is_success"`
}

func OK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{IsSuccess: true, Data: data})
}

func Created(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, Response{IsSuccess: true, Data: data})
}

func HandleError(ctx *gin.Context, err error) {
	var (
		code    int
		message string

		derr *derrors.Error
	)

	switch {
	case errors.As(err, &derr):
		code = derr.Code
		message = derr.Message
	default:
		code = derrors.ErrInternalServerError.Code
		message = derrors.ErrInternalServerError.Message
	}

	_ = ctx.Error(err)

	resp := Response{
		IsSuccess: false,
		Error:     message,
	}

	ctx.JSON(code, resp)
	ctx.Abort()
}
