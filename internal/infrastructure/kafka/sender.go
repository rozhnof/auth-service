package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type MessageSender struct {
	producer Producer
}

func NewMessageSender(producer Producer) MessageSender {
	return MessageSender{
		producer: producer,
	}
}

func (s MessageSender) SendMessage(ctx context.Context, message Message) error {
	kafkaMessage, err := s.buildKafkaMessage(message)
	if err != nil {
		return err
	}

	if err := s.producer.SendMessage(kafkaMessage); err != nil {
		return err
	}

	return nil
}

func (s MessageSender) SendMessages(ctx context.Context, messages []Message) error {
	kafkaMessageList := make([]*sarama.ProducerMessage, 0, len(messages))

	for _, message := range messages {
		kafkaMessage, err := s.buildKafkaMessage(message)
		if err != nil {
			return err
		}

		kafkaMessageList = append(kafkaMessageList, kafkaMessage)
	}

	if err := s.producer.SendMessages(kafkaMessageList); err != nil {
		return err
	}

	return nil
}

func (s MessageSender) buildKafkaMessage(message Message) (*sarama.ProducerMessage, error) {
	return &sarama.ProducerMessage{
		Key:       sarama.StringEncoder(message.Key.String()),
		Value:     sarama.ByteEncoder(message.Value),
		Topic:     message.Topic,
		Partition: -1,
	}, nil
}
