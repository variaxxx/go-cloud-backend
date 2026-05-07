package auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	RefreshTtl time.Duration `envconfig:"REFRESH_TTL" required:"true"`
}

func NewConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("AUTH", &cfg); err != nil {
		return config{}, fmt.Errorf("Auth config parse: %w", err)
	}

	return cfg, nil
}
