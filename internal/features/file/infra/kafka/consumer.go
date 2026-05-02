package file_kafka

import (
	core_logger "cloud/internal/core/logger"
	file_app "cloud/internal/features/file/application"
	infra_kafka "cloud/internal/infra/kafka"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type fileUploadedMessageHandler struct {
	handler file_app.FileUploadedEventHandler
}

type FileUploadedConsumer struct {
	consumer infra_kafka.Consumer
	handler  fileUploadedMessageHandler
	topic    string
	groupID  string
}

func NewFileUploadedConsumer(
	consumer infra_kafka.Consumer,
	handler file_app.FileUploadedEventHandler,
	topic string,
	groupID string,
) *FileUploadedConsumer {
	return &FileUploadedConsumer{
		consumer: consumer,
		handler: fileUploadedMessageHandler{
			handler: handler,
		},
		topic:   topic,
		groupID: groupID,
	}
}

func (c *FileUploadedConsumer) Run(
	ctx context.Context,
) error {
	log := core_logger.FromContext(ctx)
	log.Info(
		"File uploaded consumer started",
		zap.String("topic", c.topic),
		zap.String("group_id", c.groupID),
	)

	return c.consumer.Run(ctx, c.handler)
}

func (c *FileUploadedConsumer) Close() error {
	return c.consumer.Close()
}

func (h fileUploadedMessageHandler) Handle(
	ctx context.Context,
	message infra_kafka.Message,
) error {
	var event file_app.FileUploadedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return infra_kafka.NewNonRetryableError(
			fmt.Errorf("decode file uploaded event: %w", err),
		)
	}

	if err := h.handler.HandleUploaded(ctx, event); err != nil {
		return fmt.Errorf("handle file uploaded event: %w", err)
	}

	return nil
}
