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
	mu     sync.Mutex
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
	c.mu.Lock()

	if c.closed {
		c.mu.Unlock()
		return
	}

	c.mu.Unlock()

	go c.writePump()
	go c.readPump()
}

func (c *Client) Send(message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return derrors.ErrClientClosed
	}

	select {
	case c.send <- message:
		return nil
	default:
		return derrors.ErrSendBufferFull
	}
}

func (c *Client) Close(mustUnreg bool, code wslib.StatusCode, reason string) {
	c.mu.Lock()

	if c.closed {
		c.mu.Unlock()
		return
	}

	c.closed = true
	close(c.send)
	c.cancel()

	c.mu.Unlock()

	if mustUnreg {
		c.hub.Unregister(c)
	}

	_ = c.conn.Close(code, reason)
}
