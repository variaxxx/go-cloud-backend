package app_load

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APIBaseURL     string        `envconfig:"API_BASE_URL" default:"http://localhost:8000"`
	Username       string        `envconfig:"USERNAME" default:"load-tester"`
	Password       string        `envconfig:"PASSWORD" default:"load-tester-password"`
	UploadCount    int           `envconfig:"UPLOAD_COUNT" default:"50"`
	Concurrency    int           `envconfig:"CONCURRENCY" default:"5"`
	RequestTimeout time.Duration `envconfig:"REQUEST_TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("LOAD", &cfg); err != nil {
		return Config{}, fmt.Errorf("LOAD config parse: %w", err)
	}

	return cfg, nil
}
