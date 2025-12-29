package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/service"
)

type handler struct {
	service *service.Service
}

func New(log *zap.Logger, cfg *config.APIServer, svc *service.Service) *http.Server {
	hdl := &handler{service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.Logging(log))

	authMiddleware := middlewares.Auth(svc.AccessJWTController)

	apiRouter := router.Group("/api")
	{
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", hdl.registerUser)
			authRouter.POST("/login", hdl.loginUser)
			authRouter.POST("/logout-all", authMiddleware, hdl.logoutUserAll)

			authRouter.Use(middlewares.RefreshAuth(svc.RefreshJWTController, svc.RepoManager))
			authRouter.POST("/refresh", hdl.refreshToken)
			authRouter.POST("/logout", hdl.logoutUser)
		}

		apiRouter.Use(authMiddleware)
		usersRouter := apiRouter.Group("/users")
		{
			usersRouter.GET("/me", hdl.getCurrentUser)
			usersRouter.PATCH("/me", hdl.updateCurrentUser)
			usersRouter.DELETE("/me", hdl.deleteCurrentUser)

			usersRouter.GET("/search", hdl.searchUsers)
		}

		chatsRouter := apiRouter.Group("/chats")
		{
			chatsRouter.POST("/", hdl.createChat)
			chatsRouter.GET("/", hdl.getChats)

			chatsRouter.POST("/:chat_id/messages", hdl.sendMessage)
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
