package infra_kafka

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Brokers      []string      `envconfig:"BROKERS"`
	ClientID     string        `envconfig:"CLIENT_ID" default:"cloud-backend"`
	BatchTimeout time.Duration `envconfig:"BATCH_TIMEOUT" default:"1s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
}

func NewConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("KAFKA", &cfg); err != nil {
		return config{}, fmt.Errorf("kafka config parse: %w", err)
	}

	return cfg, nil
}
