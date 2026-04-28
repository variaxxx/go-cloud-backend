package auth_jwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret string        `envconfig:"SECRET" required:"true"`
	Ttl    time.Duration `envconfig:"TTL" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("JWT config parse: %w", err)
	}

	return config, nil
}
