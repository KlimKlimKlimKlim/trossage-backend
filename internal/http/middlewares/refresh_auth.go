package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/controller"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/jwt"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

func RefreshAuth(jwtController *jwt.Controller, rm controller.IRepoManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractBearerToken(ctx)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to extract token: %w", err))
			return
		}

		userID, err := jwtController.ProcessToken(tokenString)
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to process token: %w", err))
			return
		}

		tokenHash := hashToken(tokenString)
		token, err := rm.Repo().SelectRefreshToken(ctx, userID, tokenHash)
		if err != nil {
			if errors.Is(err, derrors.ErrTokenNotFound) {
				err = fmt.Errorf("token not found: %w", derrors.ErrUnauthorized)
			}

			response.HandleError(ctx, fmt.Errorf("failed to select token: %w", err))
			return
		}

		if !isTokenValid(token) {
			response.HandleError(ctx, fmt.Errorf("invalid token: %w", derrors.ErrUnauthorized))
			return
		}

		ctx.Set(contextKeyUserID, userID)
		ctx.Set(contextKeyTokenID, token.ID)
		ctx.Next()
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func isTokenValid(token models.Token) bool {
	now := time.Now()

	if token.ExpiresAt.Before(now) {
		return false
	}

	if !token.RevokedAt.IsZero() {
		return false
	}

	return true
}
