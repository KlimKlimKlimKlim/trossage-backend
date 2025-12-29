package postgres

import (
	"context"
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type IRepository interface { //nolint:interfacebloat // it's okay for repository interfaces
	InsertUser(ctx context.Context, user models.AuthUser) (models.AuthUser, error)
	SelectAuthUserByLogin(ctx context.Context, login string) (models.AuthUser, error)
	SelectAuthUserByID(ctx context.Context, userID int64) (models.AuthUser, error)
	SelectUserByID(ctx context.Context, userID int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.AuthUser) (models.AuthUser, error)
	DeleteUser(ctx context.Context, userID int64) error
	SelectUsersByLoginPrefix(ctx context.Context, userID int64, query string, limit, offset int) ([]models.User, error)
	CountUsersByLoginPrefix(ctx context.Context, userID int64, query string) (int, error)

	InsertRefreshToken(ctx context.Context, token models.Token) (models.Token, error)
	SelectRefreshToken(ctx context.Context, userID int64, tokenHash string) (models.Token, error)
	RevokeRefreshTokenByID(ctx context.Context, tokenID int64) error
	RevokeRefreshTokensByUserID(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context, olderThan time.Time) (int64, error)
	DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error)
}
