package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/jwt"
)

func Auth(jwtController *jwt.Controller) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extractBearerToken(ctx)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to extract token: %w", err))
			return
		}

		userID, err := jwtController.ProcessToken(token)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to process token: %w", err))
			return
		}

		ctx.Set(contextKeyUserID, userID)
		ctx.Next()
	}
}
