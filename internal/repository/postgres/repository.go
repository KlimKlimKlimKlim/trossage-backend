package postgres

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/repository"
)

type Repository struct {
	querier Querier
}

var _ repository.Repository = (*Repository)(nil)

func NewRepository(q Querier) *Repository {
	return &Repository{querier: q}
}
