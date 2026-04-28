package auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RefreshTtl time.Duration `envconfig:"REFRESH_TTL" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("AUTH", &config); err != nil {
		return Config{}, fmt.Errorf("Auth config parse: %w", err)
	}

	return config, nil
}
