package logger

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
)

func New(cfg config.Logger) (*zap.Logger, error) {
	switch cfg.Env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		return zap.NewProduction()
	case "dev":
		gin.SetMode(gin.DebugMode)
		return zap.NewDevelopment()
	case "test":
		gin.SetMode(gin.TestMode)
		return zap.NewExample(), nil
	default:
		return nil, errors.New("unknown logger environment")
	}
}
