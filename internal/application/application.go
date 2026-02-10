package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/server/api"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/server/system"
	"github.com/GlaciemArgentum/trossage-backend/internal/logger"
	"github.com/GlaciemArgentum/trossage-backend/internal/repository/postgres"
	"github.com/GlaciemArgentum/trossage-backend/internal/service"
	ws "github.com/GlaciemArgentum/trossage-backend/internal/websocket/hub"
	"github.com/GlaciemArgentum/trossage-backend/internal/worker/tokencleanup"
	"github.com/GlaciemArgentum/trossage-backend/migrations"
)

type App struct {
	log    *zap.Logger
	config *config.Config

	errorGroup    *errgroup.Group
	errorGroupCtx context.Context //nolint:containedctx // errgroup must share one context across workers

	pool *pgxpool.Pool

	service *service.Service
	wsHub   *ws.Hub

	tokenCleanupWorker *tokencleanup.Worker

	systemServer *http.Server
	apiServer    *http.Server
}

func New(cfg *config.Config) (*App, error) {
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create new logger: %w", err)
	}

	app := App{
		log:    log,
		config: cfg,
	}

	return &app, nil
}

func (a *App) Init(ctx context.Context) error {
	pool, err := postgres.New(ctx, &a.config.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	a.pool = pool
	a.log.Info("Database connected")

	if err = migrations.RunMigrations(a.config.Postgres.URL()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.log.Info("Migrations run")

	a.wsHub = ws.New(a.log, &a.config.Server.WebSocket)
	repo := postgres.NewRepository(pool)
	a.service = service.New(a.config, repo, a.wsHub)

	a.tokenCleanupWorker = tokencleanup.New(a.log, repo, &a.config.Worker.TokenCleanup)

	a.systemServer = system.New(&a.config.Server.System)
	a.apiServer = api.New(a.log, a.config, a.service)

	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.errorGroup, a.errorGroupCtx = errgroup.WithContext(ctx)

	a.errorGroup.Go(func() error {
		a.log.Info("WebSocket Hub starting")
		a.wsHub.Run(a.errorGroupCtx)

		return nil
	})

	a.errorGroup.Go(func() error {
		a.log.Info("System HTTP server starting", zap.String("address", a.systemServer.Addr))

		if err := a.systemServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("system HTTP server: %w", err)
		}

		return nil
	})

	a.errorGroup.Go(func() error {
		a.log.Info("API HTTP server starting", zap.String("address", a.apiServer.Addr))

		if err := a.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("api HTTP server: %w", err)
		}

		return nil
	})

	a.errorGroup.Go(func() error {
		if err := a.tokenCleanupWorker.Run(a.errorGroupCtx); err != nil {
			return fmt.Errorf("token cleanup worker: %w", err)
		}

		return nil
	})

	return nil
}

func (a *App) Wait() error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		err := a.errorGroup.Wait()
		errChan <- err
	}()

	select {
	case sig := <-quit:
		a.log.Info("Received shutdown signal", zap.String("signal", sig.String()))
		return nil
	case err := <-errChan:
		if err != nil {
			a.log.Error("Component error", zap.Error(err))
			return err
		}

		return nil
	}
}

func (a *App) Stop() error {
	a.log.Info("Stopping application")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	if a.apiServer != nil {
		if err := a.apiServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("API HTTP server shutdown error", zap.Error(err))
		} else {
			a.log.Info("API HTTP server stopped")
		}
	}

	if a.wsHub != nil {
		a.wsHub.Stop(shutdownCtx)
		a.log.Info("WebSocket Hub stopped")
	}

	if a.pool != nil {
		a.pool.Close()
		a.log.Info("Database disconnected")
	}

	if a.systemServer != nil {
		if err := a.systemServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("System HTTP server shutdown error", zap.Error(err))
		} else {
			a.log.Info("System HTTP server stopped")
		}
	}

	_ = a.log.Sync()

	a.log.Info("Application stopped successfully")

	return nil
}
