package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
)

func New(cfg *config.SystemServer) *http.Server {
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/health", health)

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
