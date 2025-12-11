package controller

import (
	"context"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

type IRepoManager interface {
	Repo() postgres.IRepository
	InTx(ctx context.Context, fn func(tx postgres.IRepository) error) (err error)
}
