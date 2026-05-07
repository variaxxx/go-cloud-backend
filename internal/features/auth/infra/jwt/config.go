package auth_jwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Secret string        `envconfig:"SECRET" required:"true"`
	Ttl    time.Duration `envconfig:"TTL" required:"true"`
}

func NewConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("JWT", &cfg); err != nil {
		return config{}, fmt.Errorf("JWT config parse: %w", err)
	}

	return cfg, nil
}
