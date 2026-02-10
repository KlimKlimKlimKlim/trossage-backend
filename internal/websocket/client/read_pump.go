package client

import (
	wslib "github.com/coder/websocket"

	ws "github.com/GlaciemArgentum/trossage-backend/internal/websocket"
)

// readPump reads messages from the WebSocket connection.
// Messages are intentionally discarded as this is a server->client only channel.
// The pump keeps the connection alive and detects client disconnection.
func (c *Client) readPump() {
	defer func() {
		c.Close(true, wslib.StatusGoingAway, ws.ReasonConnectionError)
	}()

	for {
		_, _, err := c.conn.Read(c.ctx)
		if err != nil {
			return
		}
	}
}
