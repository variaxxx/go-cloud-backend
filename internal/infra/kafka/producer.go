package infra_kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Header struct {
	Key   string
	Value []byte
}

type Message struct {
	Topic   string
	Key     []byte
	Value   []byte
	Headers []Header
}

type Producer interface {
	Publish(
		ctx context.Context,
		message Message,
	) error

	CreateTopic(
		ctx context.Context,
		name string,
		numPartitions int,
	) error

	Close() error
}

type KafkaProducer struct {
	writer  *kafka.Writer
	brokers []string
}

func NewProducer(cfg config) *KafkaProducer {
	return &KafkaProducer{
		brokers: cfg.Brokers,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: cfg.BatchTimeout,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			Transport: &kafka.Transport{
				ClientID: cfg.ClientID,
			},
		},
	}
}

func (p *KafkaProducer) Publish(
	ctx context.Context,
	message Message,
) error {
	if message.Topic == "" {
		return fmt.Errorf("publish kafka message: empty topic")
	}

	headers := make([]kafka.Header, 0, len(message.Headers))
	for _, header := range message.Headers {
		headers = append(headers, kafka.Header{
			Key:   header.Key,
			Value: header.Value,
		})
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic:   message.Topic,
		Key:     message.Key,
		Value:   message.Value,
		Headers: headers,
	}); err != nil {
		return fmt.Errorf("publish kafka message: %w", err)
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("close kafka writer: %w", err)
	}

	return nil
}

func (p *KafkaProducer) CreateTopic(
	ctx context.Context,
	name string,
	numPartitions int,
) error {
	if name == "" {
		return fmt.Errorf("create kafka topic: empty topic name")
	}

	if numPartitions <= 0 {
		return fmt.Errorf("create kafka topic: invalid partitions count")
	}

	if len(p.brokers) == 0 {
		return fmt.Errorf("create kafka topic: no brokers configured")
	}

	conn, err := kafka.DialContext(ctx, "tcp", p.brokers[0])
	if err != nil {
		return fmt.Errorf("create kafka topic: dial broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("create kafka topic: get controller: %w", err)
	}

	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("create kafka topic: dial controller: %w", err)
	}
	defer controllerConn.Close()

	if err := controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     numPartitions,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("create kafka topic %q: %w", name, err)
	}

	return nil
}
