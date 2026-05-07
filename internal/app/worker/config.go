package app_worker

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MetricsAddr            string        `envconfig:"METRICS_ADDR" default:":8082"`
	MetricsShutdownTimeout time.Duration `envconfig:"METRICS_SHUTDOWN_TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("WORKER", &cfg); err != nil {
		return Config{}, fmt.Errorf("WORKER config parse: %w", err)
	}

	return cfg, nil
}
