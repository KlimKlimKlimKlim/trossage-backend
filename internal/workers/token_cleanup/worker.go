package token_cleanup

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/controller"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/logger"
)

type Worker struct {
	log        *zap.Logger
	repository controller.IRepoManager
	cfg        *config.TokenCleanupConfig
}

func New(log *zap.Logger, rm controller.IRepoManager, cfg *config.TokenCleanupConfig) *Worker {
	return &Worker{
		log:        log.With(zap.String(logger.WorkerField, workerName)),
		repository: rm,
		cfg:        cfg,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("Starting worker")

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		if err := w.cleanup(ctx); err != nil {
			w.log.Error("Cleanup failed", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			w.log.Info("Stopping worker")
			return nil
		case <-ticker.C:
		}
	}
}

func (w *Worker) cleanup(ctx context.Context) error {
	w.log.Debug("Starting token cleanup")

	now := time.Now()

	expiredThreshold := now.Add(-w.cfg.ExpiredRetention)
	expiredCount, err := w.repository.Repo().DeleteExpiredTokens(ctx, expiredThreshold)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	revokedThreshold := now.Add(-w.cfg.RevokedRetention)
	revokedCount, err := w.repository.Repo().DeleteRevokedTokens(ctx, revokedThreshold)
	if err != nil {
		return fmt.Errorf("failed to delete revoked tokens: %w", err)
	}

	w.log.Debug("Token cleanup completed",
		zap.Int64("expired_deleted", expiredCount),
		zap.Int64("revoked_deleted", revokedCount),
	)

	return nil
}
