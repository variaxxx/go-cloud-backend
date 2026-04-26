package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("HTTP server config parse: %w", err)
	}

	return config, nil
}

func NewConfigRequired() Config {
	config, err := NewConfig()

	if err != nil {
		err = fmt.Errorf("HTTP server config parse: %w", err)
		panic(err)
	}

	return config
}
