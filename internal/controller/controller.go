package controller

import (
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/hasher"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/jwt"
)

type Controller struct {
	hasher *hasher.Hasher

	accessJWTController  *jwt.Controller
	refreshJWTController *jwt.Controller

	repoManager IRepoManager
}

func New(cfg *config.Config, repoManager IRepoManager) *Controller {
	return &Controller{
		hasher:               hasher.New(&cfg.Hasher),
		accessJWTController:  jwt.New(&cfg.JWT.Access, jwt.AccessType),
		refreshJWTController: jwt.New(&cfg.JWT.Refresh, jwt.RefreshType),
		repoManager:          repoManager,
	}
}
