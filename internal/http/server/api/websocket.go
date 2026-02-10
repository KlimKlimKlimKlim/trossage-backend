package api

import (
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GlaciemArgentum/trossage-backend/internal/config"
	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/jwt"
	ws "github.com/GlaciemArgentum/trossage-backend/internal/websocket/client"
)

// connectWebSocket establishes a WebSocket connection for real-time chat events.
//
//	@Summary		Connect to WebSocket
//	@Description	Establish WebSocket connection for receiving real-time chat events (server -> client only)
//	@Tags			websocket
//	@Security		BearerAuth
//	@Param			token	query		string				false	"JWT token (alternative to Authorization header for WebSocket clients)"
//	@Success		101		{string}	string				"Switching Protocols"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/ws [get]
func (h *handler) connectWebSocket(
	logger *zap.Logger,
	jwtController *jwt.Controller,
	cfg *config.WebSocket,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		status := http.StatusSwitchingProtocols

		var err error

		defer func() {
			logWebSocket(logger, ctx.Request, status, time.Since(start), err)
		}()

		token := ctx.Query("token")
		if token == "" {
			status, err = http.StatusUnauthorized, derrors.ErrUnauthorized
			ctx.Status(status)

			return
		}

		userID, err := jwtController.ProcessToken(token)
		if err != nil {
			status = http.StatusUnauthorized
			ctx.Status(status)

			return
		}

		writer := newStatusCapture(ctx.Writer)

		conn, err := websocket.Accept(writer, ctx.Request, &websocket.AcceptOptions{
			CompressionMode: websocket.CompressionContextTakeover,
			OriginPatterns:  cfg.AllowedOrigins,
		})
		if err != nil {
			status = writer.statusCode

			return
		}

		client := ws.NewClient(h.service.WSHub, conn, userID)
		if h.service.WSHub.Register(client) {
			client.Run()
		}
	}
}

// statusCapture wraps http.ResponseWriter to capture the status code.
// This is needed because websocket.Accept may write error responses directly.
type statusCapture struct {
	http.ResponseWriter
	statusCode int
}

func newStatusCapture(w http.ResponseWriter) *statusCapture {
	return &statusCapture{
		ResponseWriter: unwrapResponseWriter(w),
		statusCode:     http.StatusInternalServerError,
	}
}

func (sc *statusCapture) WriteHeader(code int) {
	sc.statusCode = code
	sc.ResponseWriter.WriteHeader(code)
}

func (sc *statusCapture) Unwrap() http.ResponseWriter {
	return sc.ResponseWriter
}

// unwrapResponseWriter extracts the underlying http.ResponseWriter from wrappers.
// This bypasses Gin's ResponseWriter to allow proper WebSocket hijacking.
func unwrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	for {
		unwrapper, ok := w.(interface{ Unwrap() http.ResponseWriter })
		if !ok {
			return w
		}

		w = unwrapper.Unwrap()
	}
}

func logWebSocket(logger *zap.Logger, req *http.Request, status int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
		zap.String("user_agent", req.UserAgent()),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	switch {
	case status >= http.StatusInternalServerError:
		logger.Error("WebSocket request", fields...)
	case status >= http.StatusBadRequest:
		logger.Warn("WebSocket request", fields...)
	default:
		logger.Info("WebSocket request", fields...)
	}
}
