package client

import (
	wslib "github.com/coder/websocket"

	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

func (c *Client) readPump() {
	defer func() {
		c.Close(true, wslib.StatusGoingAway, ws.ReasonServerShutdown)
	}()

	for {
		_, _, err := c.conn.Read(c.ctx)
		if err != nil {
			return
		}
	}
}
