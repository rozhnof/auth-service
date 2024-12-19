package outbox

import (
	"context"
	"log/slog"
	"time"

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
	logger     *slog.Logger
	tracer     trace.Tracer
}

func NewKafkaOutboxSender(txManager trm.TransactionManager, messageSender MessageSender, log *slog.Logger, tracer trace.Tracer) *KafkaOutboxSender {
	return &KafkaOutboxSender{
		repository: NewRepository(txManager, log, tracer),
		txManager:  txManager,
		sender:     messageSender,
		logger:     log,
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

func (s *KafkaOutboxSender) Run(ctx context.Context, topics []string, batchSize int32, readInterval time.Duration) error {
	ctx, span := s.tracer.Start(ctx, "OutboxMessageSender.Run")
	defer span.End()

	ticker := time.NewTicker(readInterval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		for _, topic := range topics {
			if err := s.Read(ctx, topic, batchSize); err != nil {
				s.logger.Warn("failed read from postgres and send message to kafka", slog.String("error", err.Error()))
			}
		}
	}
}
