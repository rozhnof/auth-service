package auth

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rozhnof/auth-service/internal/infrastructure/kafka"
)

type IMessageSender interface {
	SendMessage(context.Context, kafka.Message) error
	SendMessages(context.Context, []kafka.Message) error
}

type MessageSender[T any] struct {
	sender IMessageSender
	topic  string
}

func NewMessageSender[T any](sender IMessageSender, topic string) MessageSender[T] {
	return MessageSender[T]{
		sender: sender,
		topic:  topic,
	}
}

func (s MessageSender[T]) SendMessage(ctx context.Context, message T) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	outboxMessage := kafka.Message{
		Key:   uuid.New(),
		Value: bytes,
		Topic: s.topic,
	}

	if err := s.sender.SendMessage(ctx, outboxMessage); err != nil {
		return err
	}

	return nil
}

func (s MessageSender[T]) SendMessages(ctx context.Context, messages []T) error {
	outboxMessages := make([]kafka.Message, 0, len(messages))

	for _, message := range messages {
		bytes, err := json.Marshal(message)
		if err != nil {
			return err
		}

		outboxMessage := kafka.Message{
			Key:   uuid.New(),
			Value: bytes,
			Topic: s.topic,
		}

		outboxMessages = append(outboxMessages, outboxMessage)
	}

	if err := s.sender.SendMessages(ctx, outboxMessages); err != nil {
		return err
	}

	return nil
}
