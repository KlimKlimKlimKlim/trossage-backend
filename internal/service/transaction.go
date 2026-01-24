package service

import (
	"context"
	"fmt"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/repository"
)

func (s *Service) inTxWithOpts(ctx context.Context, readOnly bool, fn func(*Service) error) (err error) {
	// If already in transaction - reuse it without creating new wrapper.
	if txRepo, ok := s.Repo.(*txRepository); ok {
		if txRepo.readOnly && !readOnly {
			return derrors.ErrReadWriteTx
		}

		return fn(s)
	}

	repo, commiter, err := s.Repo.WithTx(ctx, readOnly)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			_ = commiter.Rollback(ctx)
			err = fmt.Errorf("%w: %v", derrors.ErrPanicInTx, rec)

			return
		}

		if err != nil {
			if rbErr := commiter.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("original error (rollback also failed): %w, rollback error: %w", err, rbErr)
			}

			return
		}

		if cmErr := commiter.Commit(ctx); cmErr != nil {
			err = fmt.Errorf("tx commit error: %w", cmErr)
		}
	}()

	stx := &Service{
		hasher:               s.hasher,
		AccessJWTController:  s.AccessJWTController,
		RefreshJWTController: s.RefreshJWTController,
		WSHub:                s.WSHub,
		Repo:                 newTxRepository(repo, readOnly),
	}

	err = fn(stx)

	return err
}

func (s *Service) InTx(ctx context.Context, fn func(*Service) error) (err error) {
	return s.inTxWithOpts(ctx, false, fn)
}

func (s *Service) InReadOnlyTx(ctx context.Context, fn func(*Service) error) (err error) {
	return s.inTxWithOpts(ctx, true, fn)
}

type txRepository struct {
	repository.Repository
	readOnly bool
}

func newTxRepository(r repository.Repository, readOnly bool) *txRepository {
	return &txRepository{
		Repository: r,
		readOnly:   readOnly,
	}
}
