package postgres

import (
	"context"
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type IRepository interface {
	InsertUser(ctx context.Context, user models.User) (models.User, error)
	SelectUserByLogin(ctx context.Context, login string) (models.User, error)
	SelectUserByID(ctx context.Context, userID int64) (models.User, error)

	InsertRefreshToken(ctx context.Context, token models.Token) (models.Token, error)
	SelectRefreshToken(ctx context.Context, userID int64, tokenHash string) (models.Token, error)
	RevokeRefreshTokenByID(ctx context.Context, tokenID int64) error
	RevokeRefreshTokensByUserID(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context, olderThan time.Time) (int64, error)
	DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error)
}
