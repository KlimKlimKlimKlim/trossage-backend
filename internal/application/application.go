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

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/controller"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/servers/api"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/servers/system"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/logger"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository"
	"github.com/KlimKlimKlimKlim/trossage-backend/migrations"
)

type App struct {
	log    *zap.Logger
	config *config.Config

	errorGroup    *errgroup.Group
	errorGroupCtx context.Context

	pool *pgxpool.Pool

	controller *controller.Controller

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

	if err = migrations.RunMigrations(&a.config.Postgres); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.log.Info("Migrations run")

	repoManager := postgres.NewManager(pool, repository.NewRepository(pool))
	a.controller = controller.New(a.config, repoManager)

	a.systemServer = system.New(&a.config.Server.System)
	a.apiServer = api.New(a.log, &a.config.Server.API, a.controller)

	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.errorGroup, a.errorGroupCtx = errgroup.WithContext(ctx)

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

	return nil
}

func (a *App) Wait() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		if err := a.errorGroup.Wait(); err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-quit:
		a.log.Info("Received shutdown signal", zap.String("signal", sig.String()))
		return nil
	case err := <-errChan:
		a.log.Error("Component error", zap.Error(err))
		return err
	}
}

func (a *App) Stop() error {
	a.log.Info("Stopping application")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	if a.systemServer != nil {
		if err := a.systemServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("System HTTP server shutdown error", zap.Error(err))
		} else {
			a.log.Info("System HTTP server stopped")
		}
	}

	if a.apiServer != nil {
		if err := a.apiServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("API HTTP server shutdown error", zap.Error(err))
		} else {
			a.log.Info("API HTTP server stopped")
		}
	}

	if a.pool != nil {
		a.pool.Close()
		a.log.Info("Database disconnected")
	}

	_ = a.log.Sync()

	a.log.Info("Application stopped successfully")

	return nil
}
