package hub

import (
	wslib "github.com/coder/websocket"
	"go.uber.org/zap"

	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (h *Hub) Register(client ws.IClient) bool {
	h.mu.RLock()
	stopped := h.stopped
	h.mu.RUnlock()

	if stopped {
		client.Close(false, wslib.StatusGoingAway, ws.ReasonServerShutdown)
		return false
	}

	select {
	case h.register <- client:
		return true
	default:
		client.Close(false, wslib.StatusGoingAway, ws.ReasonRegisterFailed)

		h.log.Error("Register channel full, dropping client",
			zap.Int64("user_id", client.UserID()),
		)

		return false
	}
}

func (h *Hub) Unregister(client ws.IClient) {
	h.mu.RLock()
	stopped := h.stopped
	h.mu.RUnlock()

	if stopped {
		return
	}

	select {
	case h.unregister <- client:
	default:
		h.log.Error("Unregister channel full, dropping client",
			zap.Int64("user_id", client.UserID()),
		)
	}
}
