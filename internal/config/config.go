package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   Server       `envPrefix:"SERVER_"`
	Logger   Logger       `envPrefix:"LOGGER_"`
	Postgres Postgres     `envPrefix:"POSTGRES_"`
	Hasher   Hasher       `envPrefix:"HASHER_"`
	JWT      JWT          `envPrefix:"JWT_"`
	Worker   WorkerConfig `envPrefix:"WORKER_"`
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

type Postgres struct {
	Host              string        `env:"HOST,required"`
	Port              string        `env:"PORT,required"`
	User              string        `env:"USER,required"`
	Password          string        `env:"PASSWORD,required"`
	DBName            string        `env:"DB,required"`
	SSLMode           string        `env:"SSL_MODE,required"`
	MaxConns          int32         `env:"MAX_CONNS" envDefault:"25"`
	MinConns          int32         `env:"MIN_CONNS" envDefault:"5"`
	MaxConnLifetime   time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"5m"`
	MaxConnIdleTime   time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"1m"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"15s"`
}

type Hasher struct {
	Memory     uint32 `env:"MEMORY" envDefault:"65536"`
	Iterations uint32 `env:"ITERATIONS" envDefault:"3"`
	SaltLength uint32 `env:"SALT_LENGTH" envDefault:"16"`
	KeyLength  uint32 `env:"KEY_LENGTH" envDefault:"32"`
}

type JWT struct {
	Access  JWTSetting `envPrefix:"ACCESS_"`
	Refresh JWTSetting `envPrefix:"REFRESH_"`
}

type JWTSetting struct {
	Secret   string        `env:"SECRET,required"`
	Lifetime time.Duration `env:"LIFETIME,required"`
}

type Logger struct {
	Env string `env:"ENV,required"`
}

type WorkerConfig struct {
	TokenCleanup TokenCleanupConfig `envPrefix:"TOKEN_CLEANUP_"`
}

type TokenCleanupConfig struct {
	Interval         time.Duration `env:"INTERVAL" envDefault:"15m"`
	ExpiredRetention time.Duration `env:"EXPIRED_RETENTION" envDefault:"1h"`
	RevokedRetention time.Duration `env:"REVOKED_RETENTION" envDefault:"1h"`
}

func New() (*Config, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
