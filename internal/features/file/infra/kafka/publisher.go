package file_kafka

import (
	file_app "cloud/internal/features/file/application"
	file_domain "cloud/internal/features/file/domain"
	infra_kafka "cloud/internal/infra/kafka"
	"context"
	"encoding/json"
	"fmt"
)

type fileEventPublisher struct {
	producer          infra_kafka.Producer
	topicFileUploaded string
}

func NewFileEventPublisher(
	producer infra_kafka.Producer,
	cfg config,
) (*fileEventPublisher, error) {
	if err := producer.CreateTopic(
		context.Background(),
		cfg.TopicFileUploaded,
		cfg.TopicFileUploadedPartitions,
	); err != nil {
		return nil, fmt.Errorf("create file uploaded topic: %w", err)
	}

	return &fileEventPublisher{
		producer:          producer,
		topicFileUploaded: cfg.TopicFileUploaded,
	}, nil
}

func (p *fileEventPublisher) PublishUploaded(
	ctx context.Context,
	file file_domain.File,
) error {
	event := file_app.FileUploadedEvent{
		FileID:     file.ID,
		Filename:   file.Filename,
		Mimetype:   file.Mimetype,
		StoragePath: file.StoragePath,
		UserID:     file.UserID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal file uploaded event: %w", err)
	}

	if err := p.producer.Publish(ctx, infra_kafka.Message{
		Topic: p.topicFileUploaded,
		Key:   []byte(file.ID.String()),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("publish uploaded file event: %w", err)
	}

	return nil
}
