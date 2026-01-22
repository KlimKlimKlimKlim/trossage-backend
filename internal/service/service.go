package service

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/hasher"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/jwt"
	ws "github.com/KlimKlimKlimKlim/trossage-backend/internal/websocket/hub"
)

type Service struct {
	hasher *hasher.Hasher

	AccessJWTController  *jwt.Controller
	RefreshJWTController *jwt.Controller

	RepoManager IRepoManager

	WSHub *ws.Hub
}

func New(cfg *config.Config, repoManager IRepoManager, wsHub *ws.Hub) *Service {
	return &Service{
		hasher:               hasher.New(&cfg.Hasher),
		AccessJWTController:  jwt.New(&cfg.JWT.Access, jwt.AccessType),
		RefreshJWTController: jwt.New(&cfg.JWT.Refresh, jwt.RefreshType),
		RepoManager:          repoManager,
		WSHub:                wsHub,
	}
}
