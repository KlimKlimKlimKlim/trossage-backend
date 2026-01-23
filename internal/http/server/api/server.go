package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/service"
)

type handler struct {
	service *service.Service
}

func New(log *zap.Logger, cfg *config.Config, svc *service.Service) *http.Server {
	hdl := &handler{service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logging(log))

	authMiddleware := middleware.Auth(svc.AccessJWTController)
	wsAuthMiddleware := middleware.WebSocketAuth(svc.AccessJWTController)
	rateLimitMiddleware := middleware.RateLimit(&cfg.Server.RateLimit)

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/ws", wsAuthMiddleware, hdl.connectWebSocket(&cfg.Server.WebSocket))

		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", rateLimitMiddleware, hdl.registerUser)
			authRouter.POST("/login", rateLimitMiddleware, hdl.loginUser)
			authRouter.POST("/logout-all", authMiddleware, hdl.logoutUserAll)

			authRouter.Use(middleware.RefreshAuth(svc.RefreshJWTController, svc.RepoManager))
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
			chatsRouter.GET("/:chat_id/messages", hdl.getMessages)

			chatsRouter.POST("/:chat_id/typing", hdl.sendTyping)
		}
	}

	return &http.Server{
		Addr:              cfg.Server.API.Port,
		Handler:           router,
		ReadTimeout:       cfg.Server.API.ReadTimeout,
		WriteTimeout:      cfg.Server.API.WriteTimeout,
		IdleTimeout:       cfg.Server.API.IdleTimeout,
		ReadHeaderTimeout: cfg.Server.API.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.Server.API.MaxHeaderBytes,
	}
}
