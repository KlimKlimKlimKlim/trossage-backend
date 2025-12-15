package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/controller"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
)

type state struct {
	controller *controller.Controller
}

func New(log *zap.Logger, cfg *config.APIServer, c *controller.Controller) *http.Server {
	s := &state{controller: c}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.Logging(log))

	authMiddleware := middlewares.Auth(c.AccessJWTController)

	apiRouter := router.Group("/api")
	{
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", s.registerUser)
			authRouter.POST("/login", s.loginUser)
			authRouter.POST("/logout-all", authMiddleware, s.logoutUserAll)

			authRouter.Use(middlewares.RefreshAuth(c.RefreshJWTController, c.RepoManager))
			authRouter.POST("/refresh", s.refreshToken)
			authRouter.POST("/logout", s.logoutUser)
		}

		apiRouter.Use(authMiddleware)
		usersRouter := apiRouter.Group("/users")
		{
			usersRouter.GET("/me", s.getCurrentUser)
			usersRouter.PATCH("/me", s.updateCurrentUser)
			usersRouter.DELETE("/me", s.deleteCurrentUser)
		}
	}

	return &http.Server{
		Addr:              cfg.Port,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}
}
