package file_kafka

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	TopicFileUploaded           string `envconfig:"TOPIC_FILE_UPLOADED" default:"file.uploaded"`
	TopicFileUploadedPartitions int    `envconfig:"TOPIC_FILE_UPLOADED_PARTITIONS" default:"1"`
	GroupFileUploaded           string `envconfig:"GROUP_FILE_UPLOADED" default:"file-uploaded-consumer"`
}

func NewConfig() (config, error) {
	var cfg config

	if err := envconfig.Process("FILE_KAFKA", &cfg); err != nil {
		return config{}, fmt.Errorf("file kafka config parse: %w", err)
	}

	return cfg, nil
}
