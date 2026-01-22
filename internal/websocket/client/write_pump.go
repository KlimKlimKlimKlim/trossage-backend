package client

import (
	"context"
	"time"

	wslib "github.com/coder/websocket"

	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (c *Client) writePump() {
	ticker := time.NewTicker(c.hub.PingInterval())

	defer func() {
		ticker.Stop()
		c.Close(true, wslib.StatusGoingAway, ws.ReasonServerShutdown)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			if err := c.writeMessage(message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.writePing(); err != nil {
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) writeMessage(message []byte) error {
	writeCtx, cancel := context.WithTimeout(c.ctx, c.hub.WriteTimeout())
	defer cancel()

	return c.conn.Write(writeCtx, wslib.MessageText, message)
}

func (c *Client) writePing() error {
	pingCtx, cancel := context.WithTimeout(c.ctx, c.hub.WriteTimeout())
	defer cancel()

	return c.conn.Ping(pingCtx)
}
