package service

import (
	"github.com/GlaciemArgentum/trossage-backend/internal/config"
	"github.com/GlaciemArgentum/trossage-backend/internal/hasher"
	"github.com/GlaciemArgentum/trossage-backend/internal/jwt"
	"github.com/GlaciemArgentum/trossage-backend/internal/repository"
	ws "github.com/GlaciemArgentum/trossage-backend/internal/websocket/hub"
)

type Service struct {
	hasher *hasher.Hasher

	AccessJWTController  *jwt.Controller
	RefreshJWTController *jwt.Controller

	Repo repository.Repository

	WSHub *ws.Hub
}

func New(cfg *config.Config, repo repository.Repository, wsHub *ws.Hub) *Service {
	return &Service{
		hasher:               hasher.New(&cfg.Hasher),
		AccessJWTController:  jwt.New(&cfg.JWT.Access, jwt.AccessType),
		RefreshJWTController: jwt.New(&cfg.JWT.Refresh, jwt.RefreshType),
		Repo:                 repo,
		WSHub:                wsHub,
	}
}
