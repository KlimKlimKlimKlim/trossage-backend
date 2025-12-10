package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func loggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("System HTTP request",
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}
