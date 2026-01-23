package hub

import (
	"context"

	wslib "github.com/coder/websocket"
	"go.uber.org/zap"

	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case <-ctx.Done():
			h.mu.Lock()
			h.stopped = true
			h.mu.Unlock()

			return
		}
	}
}

func (h *Hub) registerClient(client ws.IClient) {
	h.mu.Lock()

	currentCount := len(h.clients[client.UserID()])
	if h.config.MaxConnectionsPerUser > 0 && currentCount >= h.config.MaxConnectionsPerUser {
		h.mu.Unlock()

		client.Close(false, wslib.StatusGoingAway, ws.ReasonMaxConnectionsLimit)

		return
	}

	h.clients[client.UserID()] = append(h.clients[client.UserID()], client)
	count := len(h.clients[client.UserID()])
	h.mu.Unlock()

	h.log.Debug("Client connected",
		zap.Int64("user_id", client.UserID()),
		zap.Int("connections", count),
	)
}

func (h *Hub) unregisterClient(client ws.IClient) {
	h.mu.Lock()

	clients := h.clients[client.UserID()]
	for i, c := range clients {
		if c == client {
			h.clients[client.UserID()] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(h.clients[client.UserID()]) == 0 {
		delete(h.clients, client.UserID())
	}

	remaining := len(h.clients[client.UserID()])
	h.mu.Unlock()

	h.log.Debug("Client disconnected",
		zap.Int64("user_id", client.UserID()),
		zap.Int("remaining", remaining),
	)
}
