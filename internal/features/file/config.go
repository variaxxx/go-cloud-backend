package file

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	StorageDir string `envconfig:"STORAGE_DIR" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("FILE", &config); err != nil {
		return Config{}, fmt.Errorf("file config parse: %w", err)
	}

	return config, nil
}
