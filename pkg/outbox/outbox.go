package outbox

import (
	"context"
	"log/slog"

	"github.com/rozhnof/auth-service/internal/infrastructure/kafka"
	trm "github.com/rozhnof/auth-service/pkg/transaction_manager"
	"go.opentelemetry.io/otel/trace"
)

type MessageSender interface {
	SendMessage(context.Context, kafka.Message) error
	SendMessages(context.Context, []kafka.Message) error
}

type KafkaOutboxSender struct {
	repository *Repository
	txManager  trm.TransactionManager
	sender     MessageSender
	log        *slog.Logger
	tracer     trace.Tracer
}

func NewKafkaOutboxSender(txManager trm.TransactionManager, messageSender MessageSender, log *slog.Logger, tracer trace.Tracer) *KafkaOutboxSender {
	return &KafkaOutboxSender{
		repository: NewRepository(txManager, log, tracer),
		txManager:  txManager,
		sender:     messageSender,
		log:        log,
		tracer:     tracer,
	}
}

func (s *KafkaOutboxSender) SendMessage(ctx context.Context, message kafka.Message) error {
	ctx, span := s.tracer.Start(ctx, "OutboxMessageSender.SendMessage")
	defer span.End()

	if err := s.repository.Create(ctx, message); err != nil {
		return err
	}

	return nil
}

func (s *KafkaOutboxSender) SendMessages(ctx context.Context, messages []kafka.Message) error {
	ctx, span := s.tracer.Start(ctx, "OutboxMessageSender.Create")
	defer span.End()

	if err := s.repository.CreateList(ctx, messages); err != nil {
		return err
	}

	return nil
}

func (s *KafkaOutboxSender) Read(ctx context.Context, topic string, limit int32) error {
	ctx, span := s.tracer.Start(ctx, "OutboxMessageSender.Read")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		messages, err := s.repository.Read(ctx, topic, limit)
		if err != nil {
			return err
		}

		if err := s.sender.SendMessages(ctx, messages); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
