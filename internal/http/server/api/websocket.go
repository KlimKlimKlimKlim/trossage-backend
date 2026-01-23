package api

import (
	"fmt"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket/client"
)

// connectWebSocket establishes a WebSocket connection for real-time chat events
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
func (h *handler) connectWebSocket(cfg *config.WebSocket) func(*gin.Context) {
	return func(ctx *gin.Context) {
		userID, ok := middleware.GetUserID(ctx)
		if !ok {
			response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
			return
		}

		conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
			CompressionMode: websocket.CompressionContextTakeover,
			OriginPatterns:  cfg.AllowedOrigins,
		})
		if err != nil {
			response.HandleError(ctx, fmt.Errorf("failed to accept websocket: %w", err))
			return
		}

		client := ws.NewClient(h.service.WSHub, conn, userID)

		if h.service.WSHub.Register(client) {
			client.Run()
		}
	}
}
