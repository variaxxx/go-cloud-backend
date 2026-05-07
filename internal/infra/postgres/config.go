package infra_postgres

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" required:"true"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Db       string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("PG", &cfg); err != nil {
		return Config{}, fmt.Errorf("PostgreSQL config parse: %w", err)
	}

	return cfg, nil
}
