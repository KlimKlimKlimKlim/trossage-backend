package hub

import (
	"context"
	"sync"
	"time"

	wslib "github.com/coder/websocket"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/logger"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket"
)

type Hub struct {
	log    *zap.Logger
	config *config.WebSocket

	clients map[int64][]ws.IClient

	register   chan ws.IClient
	unregister chan ws.IClient

	mu      sync.RWMutex
	stopped bool
}

var _ ws.IHub = (*Hub)(nil)

func New(log *zap.Logger, cfg *config.WebSocket) *Hub {
	return &Hub{
		log:        log.With(zap.String(logger.WorkerField, workerName)),
		config:     cfg,
		clients:    make(map[int64][]ws.IClient),
		register:   make(chan ws.IClient, cfg.ChannelBufferSize),
		unregister: make(chan ws.IClient, cfg.ChannelBufferSize),
	}
}

func (h *Hub) ClientBufferSize() int {
	return h.config.ClientBufferSize
}

func (h *Hub) PingInterval() time.Duration {
	return h.config.PingInterval
}

func (h *Hub) WriteTimeout() time.Duration {
	return h.config.WriteTimeout
}

func (h *Hub) Stop(ctx context.Context) {
	h.mu.Lock()
	h.stopped = true
	h.mu.Unlock()

	h.closeAllClients(ctx)
}

func (h *Hub) closeAllClients(ctx context.Context) {
	h.mu.Lock()

	totalClients := 0
	for _, clients := range h.clients {
		totalClients += len(clients)
	}

	clientsCopy := make([]ws.IClient, 0, totalClients)
	for _, clients := range h.clients {
		clientsCopy = append(clientsCopy, clients...)
	}

	h.clients = make(map[int64][]ws.IClient)
	h.mu.Unlock()

	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(closeClientsConcurrencyLimit)

	for _, client := range clientsCopy {
		eg.Go(func() error {
			client.Close(false, wslib.StatusGoingAway, ws.ReasonServerShutdown)
			return nil
		})
	}

	_ = eg.Wait()
}
