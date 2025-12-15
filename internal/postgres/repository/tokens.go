package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/queries/tokens"
)

func (r *Repository) InsertRefreshToken(ctx context.Context, token models.Token) (models.Token, error) {
	err := r.db.QueryRow(ctx, tokens.InsertTokenQuery, token.UserID, token.TokenHash, token.ExpiresAt).Scan(
		&token.ID,
		&token.CreatedAt,
	)

	if err != nil {
		return models.Token{}, err
	}

	return token, nil
}

func (r *Repository) SelectRefreshToken(ctx context.Context, userID int64, tokenHash string) (models.Token, error) {
	var token models.Token

	err := r.db.QueryRow(ctx, tokens.SelectTokenQuery, userID, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return models.Token{}, derrors.ErrTokenNotFound
		}

		return models.Token{}, err
	}

	return token, nil
}

func (r *Repository) RevokeRefreshTokenByID(ctx context.Context, tokenID int64) error {
	result, err := r.db.Exec(ctx, tokens.RevokeTokenByIDQuery, tokenID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return derrors.ErrTokenNotFound
	}

	return nil
}

func (r *Repository) RevokeRefreshTokensByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, tokens.RevokeTokensByUserID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteExpiredTokens(ctx context.Context, olderThan time.Time) (int64, error) {
	result, err := r.db.Exec(ctx, tokens.DeleteExpiredTokensQuery, olderThan)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *Repository) DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error) {
	result, err := r.db.Exec(ctx, tokens.DeleteRevokedTokensQuery, olderThan)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
