package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
)

func extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader(authorizationHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", derrors.ErrUnauthorized
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", derrors.ErrUnauthorized
	}

	return token, nil
}

func GetUserID(ctx *gin.Context) (int64, bool) {
	userID, exists := ctx.Get(contextKeyUserID)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}

func GetTokenID(ctx *gin.Context) (int64, bool) {
	userID, exists := ctx.Get(contextKeyTokenID)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}
