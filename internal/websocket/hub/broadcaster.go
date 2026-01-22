package hub

import (
	"encoding/json"
	"errors"

	wslib "github.com/coder/websocket"
	"go.uber.org/zap"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (h *Hub) BroadcastToUser(userID int64, event *Event) {
	data, err := json.Marshal(event)
	if err != nil {
		h.log.Error("Failed to marshal event", zap.Int64("user_id", userID), zap.Error(err))
		return
	}

	h.mu.RLock()
	clients := make([]ws.IClient, len(h.clients[userID]))
	copy(clients, h.clients[userID])
	h.mu.RUnlock()

	for _, client := range clients {
		if err = client.Send(data); err != nil {
			if errors.Is(err, derrors.ErrClientClosed) {
				h.Unregister(client)
				continue
			}

			h.log.Warn("Failed to send event to client", zap.Int64("user_id", userID), zap.Error(err))
		}
	}
}

func (h *Hub) BroadcastToUsers(userIDs []int64, event *Event) {
	for _, userID := range userIDs {
		h.BroadcastToUser(userID, event)
	}
}

func (h *Hub) DisconnectUser(userID int64) {
	h.BroadcastToUser(userID, NewUserLogoutEvent())

	h.mu.RLock()
	clients := make([]ws.IClient, len(h.clients[userID]))
	copy(clients, h.clients[userID])
	h.mu.RUnlock()

	for _, client := range clients {
		client.Close(true, wslib.StatusNormalClosure, ws.ReasonUserLogout)
	}
}
