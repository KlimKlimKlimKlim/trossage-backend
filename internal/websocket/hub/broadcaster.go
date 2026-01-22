package hub

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (h *Hub) BroadcastToUser(userID int64, eventType EventType, data any) error {
	event, err := NewServerEvent(eventType, data)
	if err != nil {
		return fmt.Errorf("failed to create server event: %w", err)
	}

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	h.mu.RLock()
	clients := make([]ws.IClient, len(h.clients[userID]))
	copy(clients, h.clients[userID])
	h.mu.RUnlock()

	for _, client := range clients {
		if err = client.Send(message); err != nil {
			if errors.Is(err, derrors.ErrClientClosed) {
				h.Unregister(client)
				continue
			}

			h.log.Warn("Failed to send event to client", zap.Int64("user_id", userID), zap.Error(err))
		}
	}

	return nil
}
