package outbox

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	Key   uuid.UUID
	Value []byte
	Topic string
}

type SenderRepository interface {
	Create(ctx context.Context, msg Message) error
}

type MessageSender struct {
	repo  SenderRepository
	topic string
}

func NewMessageSender(repo SenderRepository, topic string) MessageSender {
	return MessageSender{
		repo:  repo,
		topic: topic,
	}
}

func (s MessageSender) SendMessage(ctx context.Context, msg any) error {
	outboxMsg, err := s.buildMessage(msg)
	if err != nil {
		return err
	}

	if err := s.repo.Create(ctx, outboxMsg); err != nil {
		return err
	}

	return nil
}

func (s MessageSender) buildMessage(value any) (Message, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Key:   uuid.New(),
		Value: data,
		Topic: s.topic,
	}, nil
}
