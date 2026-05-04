package app_api

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MetricsAddr            string        `envconfig:"METRICS_ADDR" default:":8081"`
	MetricsShutdownTimeout time.Duration `envconfig:"METRICS_SHUTDOWN_TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("API", &cfg); err != nil {
		return Config{}, fmt.Errorf("API config parse: %w", err)
	}

	return cfg, nil
}
