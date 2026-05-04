package file

import (
	file_app "cloud/internal/features/file/application"
	file_kafka "cloud/internal/features/file/infra/kafka"
	file_postgres "cloud/internal/features/file/infra/postgres"
	file_storage "cloud/internal/features/file/infra/storage"
	infra_kafka "cloud/internal/infra/kafka"
	infra_postgres "cloud/internal/infra/postgres"
	obs_prometheus "cloud/internal/observability/prometheus"
	"context"
)

type Worker interface {
	Run(ctx context.Context) error
	Close() error
}

type WorkerDeps struct {
	DB            infra_postgres.Pool
	Observability *obs_prometheus.Observability
}

func NewUploadedWorker(
	deps WorkerDeps,
) (*file_kafka.FileUploadedConsumer, error) {
	fileConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}

	kafkaConfig, err := infra_kafka.NewConfig()
	if err != nil {
		return nil, err
	}

	fileKafkaConfig, err := file_kafka.NewConfig()
	if err != nil {
		return nil, err
	}

	fileRepo := file_postgres.NewFileRepository(deps.DB)
	storage := file_storage.NewLocalFileStorage(fileConfig.StorageDir)
	metrics := deps.Observability.File

	handler := file_app.NewFileUploadedEventHandler(fileRepo, storage, metrics)
	consumer := infra_kafka.NewConsumer(kafkaConfig, infra_kafka.ConsumerParams{
		Topic:   fileKafkaConfig.TopicFileUploaded,
		GroupID: fileKafkaConfig.GroupFileUploaded,
	})

	return file_kafka.NewFileUploadedConsumer(
		consumer,
		handler,
		fileKafkaConfig.TopicFileUploaded,
		fileKafkaConfig.GroupFileUploaded,
	), nil
}
