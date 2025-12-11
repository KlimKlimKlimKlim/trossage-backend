package postgres

import (
	"context"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
)

type IRepository interface {
	InsertUser(ctx context.Context, user models.User) (models.User, error)
	SelectUserByLogin(ctx context.Context, login string) (models.User, error)

	InsertToken(ctx context.Context, token models.Token) (models.Token, error)
}
