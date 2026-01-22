package websocket

import (
	"time"

	"github.com/coder/websocket"
)

type IHub interface { //nolint:iface // used by external packages
	ClientBufferSize() int
	PingInterval() time.Duration
	WriteTimeout() time.Duration
	Unregister(client IClient)
}

type IClient interface {
	Close(graceful bool, code websocket.StatusCode, reason string)
	Send(message []byte) error
	UserID() int64
}
