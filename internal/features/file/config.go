package file

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	StorageDir string `envconfig:"STORAGE_DIR" required:"true"`
}

func NewConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("FILE", &cfg); err != nil {
		return config{}, fmt.Errorf("file config parse: %w", err)
	}

	return cfg, nil
}
