package client

import (
	"context"
	"sync"

	wslib "github.com/coder/websocket"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

type Client struct {
	ctx    context.Context //nolint:containedctx // context lifetime equals client lifetime
	cancel context.CancelFunc
	hub    ws.IHub
	conn   *wslib.Conn
	send   chan []byte
	userID int64
	mu     sync.RWMutex
	closed bool
}

var _ ws.IClient = (*Client)(nil)

func NewClient(hub ws.IHub, conn *wslib.Conn, userID int64) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, hub.ClientBufferSize()),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) UserID() int64 {
	return c.userID
}

func (c *Client) Run() {
	c.mu.RLock()

	if c.closed {
		c.mu.RUnlock()
		return
	}

	c.mu.RUnlock()

	go c.writePump()
	go c.readPump()
}

func (c *Client) Send(event []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return derrors.ErrClientClosed
	}

	select {
	case c.send <- event:
		return nil
	default:
		return derrors.ErrSendBufferFull
	}
}

func (c *Client) Close(graceful bool, code wslib.StatusCode, reason string) {
	c.mu.Lock()

	if c.closed {
		c.mu.Unlock()
		return
	}

	c.closed = true
	close(c.send)
	c.cancel()

	c.mu.Unlock()

	if graceful {
		c.hub.Unregister(c)
	}

	_ = c.conn.Close(code, reason)
}
