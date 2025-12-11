package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
)

type Response struct {
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error"`
	Data      any    `json:"data"`
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{IsSuccess: true, Data: data})
}

func HandleError(c *gin.Context, err error) {
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
		code = http.StatusInternalServerError
		message = "internal error"
	}

	_ = c.Error(err)

	resp := Response{
		IsSuccess: false,
		Error:     message,
	}

	c.JSON(code, resp)
}
