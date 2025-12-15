package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Next()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.Error(c.Errors[0].Err))
		}

		switch {
		case c.Writer.Status() >= 500:
			logger.Error("HTTP request", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("HTTP request", fields...)
		default:
			logger.Info("HTTP request", fields...)
		}
	}
}
