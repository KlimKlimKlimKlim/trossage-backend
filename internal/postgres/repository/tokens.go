package repository

import (
	"context"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/queries/tokens"
)

func (r *Repository) InsertToken(ctx context.Context, token models.Token) (models.Token, error) {
	err := r.db.QueryRow(ctx, tokens.InsertTokenQuery, token.UserID, token.TokenHash, token.ExpiresAt).Scan(
		&token.ID,
		&token.CreatedAt,
	)

	if err != nil {
		return models.Token{}, err
	}

	return token, nil
}
