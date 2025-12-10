package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server `envPrefix:"SERVER_"`
	Logger Logger `envPrefix:"LOGGER_"`
}

type Server struct {
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	System          SystemServer  `envPrefix:"SYSTEM_"`
	API             APIServer     `envPrefix:"API_"`
}

type SystemServer struct {
	Port              string        `env:"PORT" envDefault:":8081"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"2s"`
	MaxHeaderBytes    int           `env:"MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
}

type APIServer struct {
	Port              string        `env:"PORT" envDefault:":8080"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"3s"`
	MaxHeaderBytes    int           `env:"MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
}

type Logger struct {
	Env string `env:"ENV,required"`
}

func New() (*Config, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
