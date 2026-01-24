package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/repository"
)

//nolint:ireturn // necessary for compatibility with repository.Repository
func (r *Repository) WithTx(ctx context.Context, readOnly bool) (repository.Repository, repository.Commiter, error) {
	txOptions := pgx.TxOptions{AccessMode: pgx.ReadWrite}
	if readOnly {
		txOptions.AccessMode = pgx.ReadOnly
	}

	// Type assert to Beginner interface - only *pgxpool.Pool and *pgx.Conn support it.
	beginner, ok := r.querier.(Beginner)
	if !ok {
		return nil, nil, derrors.ErrTxNotSupported
	}

	tx, err := beginner.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return NewRepository(tx), tx, nil
}
