package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env         Env `env:"ENV" env-default:"local"`
	Port        int `env:"PORT" env-default:"8080"`
	WorkerCount int `env:"WORKER_COUNT" env-default:"4"`
}

func New() (*Config, error) {
	c := &Config{}

	if err := cleanenv.ReadEnv(c); err != nil {
		return nil, err
	}

	if c.WorkerCount < 1 {
		return nil, errors.New("worker count must be greater than 0")
	}

	if c.Port < 1 || c.Port > 65535 {
		return nil, errors.New("port must be between 1 and 65535")
	}

	c.setupLogger()

	slog.Debug("config parsed", slog.Any("config", c))

	return c, nil
}

func (c *Config) setupLogger() {
	var logger *slog.Logger

	switch c.Env {
	case EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	slog.SetDefault(logger)
}
