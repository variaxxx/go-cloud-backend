package infra_kafka

import (
	core_logger "cloud/internal/core/logger"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ConsumerHandler interface {
	Handle(
		ctx context.Context,
		message Message,
	) error
}

type Consumer interface {
	Run(
		ctx context.Context,
		handler ConsumerHandler,
	) error

	Close() error
}

type ConsumerParams struct {
	Topic   string
	GroupID string
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(
	cfg config,
	params ConsumerParams,
) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			GroupID:  params.GroupID,
			Topic:    params.Topic,
			MinBytes: 1,
			MaxBytes: 10e6,
			MaxWait:  cfg.MaxWait,
		}),
	}
}

func (c *KafkaConsumer) Run(
	ctx context.Context,
	handler ConsumerHandler,
) error {
	for {
		log := core_logger.FromContext(ctx)

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("fetch kafka message: %w", err)
		}

		log.Debug(
			"Received new message",
			zap.String("topic", msg.Topic),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
		)

		message := Message{
			Topic: msg.Topic,
			Key:   msg.Key,
			Value: msg.Value,
		}

		if len(msg.Headers) > 0 {
			message.Headers = make([]Header, 0, len(msg.Headers))
			for _, header := range msg.Headers {
				message.Headers = append(message.Headers, Header{
					Key:   header.Key,
					Value: header.Value,
				})
			}
		}

		if err := handler.Handle(ctx, message); err != nil {
			fields := []zap.Field{
				zap.String("topic", msg.Topic),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
				zap.Error(err),
			}

			if IsNonRetryableError(err) {
				log.Warn("Skipping non-retryable message", fields...)

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					return fmt.Errorf("commit non-retryable kafka message: %w", err)
				}

				continue
			}

			log.Error("Failed to process message", fields...)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit kafka message: %w", err)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("close kafka reader: %w", err)
	}

	return nil
}
