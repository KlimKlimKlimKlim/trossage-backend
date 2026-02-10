package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/response"
	"github.com/GlaciemArgentum/trossage-backend/internal/jwt"
	"github.com/GlaciemArgentum/trossage-backend/internal/model"
	"github.com/GlaciemArgentum/trossage-backend/internal/repository"
)

func RefreshAuth(jwtController *jwt.Controller, rm repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractHeaderToken(ctx)
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

		token, err := rm.SelectRefreshToken(ctx, userID, tokenHash)
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

		SetUserID(ctx, userID)
		ctx.Set(contextKeyTokenID, token.ID)
		ctx.Next()
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func isTokenValid(token model.Token) bool {
	now := time.Now()

	if token.ExpiresAt.Before(now) {
		return false
	}

	if token.IsRevoked() {
		return false
	}

	return true
}
