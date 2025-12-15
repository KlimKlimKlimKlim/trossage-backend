package main

// @title           Trossage API
// @version         1.0
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
import (
	"context"
	"log"
	"os"

	_ "github.com/KlimKlimKlimKlim/trossage-backend/docs"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/application"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.New()
	if err != nil {
		log.Println("Failed to create config:", err)
		return 1
	}

	app, err := application.New(cfg)
	if err != nil {
		log.Println("Failed to create app:", err)
		return 1
	}

	ctx := context.Background()

	if err = app.Init(ctx); err != nil {
		log.Println("Failed to init app:", err)
		return 1
	}

	defer func() {
		if err = app.Stop(); err != nil {
			log.Println("Failed to stop app:", err)
		}
	}()

	if err = app.Start(ctx); err != nil {
		log.Println("Failed to start app:", err)
		return 1
	}

	if err = app.Wait(); err != nil {
		log.Println("App error:", err)
		return 1
	}

	return 0
}
