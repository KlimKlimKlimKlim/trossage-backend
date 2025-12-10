package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDHeader = "X-Request-ID"
)

func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Next()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		switch {
		case c.Writer.Status() >= 500:
			logger.Error("API HTTP request", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("API HTTP request", fields...)
		default:
			logger.Info("API HTTP request", fields...)
		}

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("Request error",
					zap.String("request_id", requestID),
					zap.Error(err),
				)
			}
		}
	}
}
