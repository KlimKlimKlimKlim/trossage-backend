package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		requestID := ctx.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Next()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", ctx.Request.UserAgent()),
		}

		if len(ctx.Errors) > 0 {
			fields = append(fields, zap.Error(ctx.Errors[0].Err))
		}

		switch {
		case ctx.Writer.Status() >= http.StatusInternalServerError:
			logger.Error("HTTP request", fields...)
		case ctx.Writer.Status() >= http.StatusBadRequest:
			logger.Warn("HTTP request", fields...)
		default:
			logger.Info("HTTP request", fields...)
		}
	}
}
