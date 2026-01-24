package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	query "github.com/KlimKlimKlimKlim/trossage-backend/internal/repository/postgres/query/token"
)

func (r *Repository) InsertRefreshToken(ctx context.Context, token model.Token) (model.Token, error) {
	err := r.querier.QueryRow(ctx, query.InsertTokenQuery, token.UserID, token.TokenHash, token.ExpiresAt).Scan(
		&token.ID,
		&token.CreatedAt,
	)
	if err != nil {
		return model.Token{}, err
	}

	return token, nil
}

func (r *Repository) SelectRefreshToken(ctx context.Context, userID int64, tokenHash string) (model.Token, error) {
	var token model.Token

	err := r.querier.QueryRow(ctx, query.SelectTokenQuery, userID, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Token{}, derrors.ErrTokenNotFound
		}

		return model.Token{}, err
	}

	return token, nil
}

func (r *Repository) RevokeRefreshTokenByID(ctx context.Context, tokenID int64) error {
	result, err := r.querier.Exec(ctx, query.RevokeTokenByIDQuery, tokenID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return derrors.ErrTokenNotFound
	}

	return nil
}

func (r *Repository) RevokeRefreshTokensByUserID(ctx context.Context, userID int64) error {
	_, err := r.querier.Exec(ctx, query.RevokeTokensByUserID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteExpiredTokens(ctx context.Context, olderThan time.Time) (int64, error) {
	result, err := r.querier.Exec(ctx, query.DeleteExpiredTokensQuery, olderThan)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *Repository) DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error) {
	result, err := r.querier.Exec(ctx, query.DeleteRevokedTokensQuery, olderThan)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
