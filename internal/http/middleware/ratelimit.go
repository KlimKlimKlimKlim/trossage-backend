package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
)

func RateLimit(cfg *config.RateLimit) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: cfg.Period,
		Limit:  cfg.Limit,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance)
}
