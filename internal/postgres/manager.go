package postgres

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository"
)

type Manager struct {
	repo IRepository
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool, repo IRepository) *Manager {
	return &Manager{
		pool: pool,
		repo: repo,
	}
}

func (m *Manager) Repo() IRepository { //nolint:ireturn // It must return interface, manager cannot return a structure.
	return m.repo
}

func (m *Manager) InTx(ctx context.Context, fn func(tx IRepository) error) (err error) {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("%w: %v\nstack: %s", derrors.ErrPanicInTx, r, debug.Stack())

			return
		}

		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("tx rollback error: %w (original error: %w)", rbErr, err)
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				err = fmt.Errorf("tx commit error: %w", cmErr)
			}
		}
	}()

	repo := repository.NewRepository(tx)
	err = fn(repo)

	return err
}
