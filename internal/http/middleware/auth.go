package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/GlaciemArgentum/trossage-backend/internal/http/response"
	"github.com/GlaciemArgentum/trossage-backend/internal/jwt"
)

func Auth(jwtController *jwt.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extractHeaderToken(ctx)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to extract token: %w", err))
			return
		}

		userID, err := jwtController.ProcessToken(token)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to process token: %w", err))
			return
		}

		SetUserID(ctx, userID)
		ctx.Next()
	}
}
