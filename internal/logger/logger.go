package logger

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
)

func New(cfg config.Logger) (*zap.Logger, error) {
	log.SetOutput(io.Discard)

	switch cfg.Env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)

		zapCfg := zap.NewProductionConfig()
		zapCfg.DisableStacktrace = true
		zapCfg.DisableCaller = true

		return zapCfg.Build()

	case "dev":
		gin.SetMode(gin.ReleaseMode)

		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.DisableStacktrace = true
		zapCfg.DisableCaller = true
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return zapCfg.Build()

	case "test":
		gin.SetMode(gin.TestMode)
		return zap.NewExample(), nil
	}

	return nil, derrors.ErrUnknownLoggerEnv
}
