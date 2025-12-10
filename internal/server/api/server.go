package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/server/api/handlers"
)

func New(log *zap.Logger, cfg *config.APIServer) *http.Server {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(log))

	router.GET("/ping", handlers.Ping)

	server := http.Server{
		Addr:              cfg.Port,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	return &server
}
