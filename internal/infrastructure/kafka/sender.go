package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type MessageSender struct {
	producer Producer
	topic    string
}

func NewMessageSender(producer Producer, topic string) MessageSender {
	return MessageSender{
		producer: producer,
		topic:    topic,
	}
}

func (s MessageSender) SendMessage(ctx context.Context, msg any) error {
	kafkaMessage, err := s.buildMessage(msg)
	if err != nil {
		return err
	}

	if err := s.producer.SendMessage(kafkaMessage); err != nil {
		return err
	}

	return nil
}

func (s MessageSender) buildMessage(value any) (*sarama.ProducerMessage, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Key:       sarama.StringEncoder(uuid.NewString()),
		Value:     sarama.ByteEncoder(data),
		Topic:     s.topic,
		Partition: -1,
	}, nil
}
