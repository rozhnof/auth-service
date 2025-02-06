package pgrepo

import (
	"context"
	"log/slog"

	db_queries "github.com/rozhnof/auth-service/internal/infrastructure/repository/queries"
	"github.com/rozhnof/auth-service/internal/pkg/outbox"
	trm "github.com/rozhnof/auth-service/pkg/transaction_manager"
	"go.opentelemetry.io/otel/trace"
)

type OutboxRepository struct {
	txManager trm.TransactionManager
	logger    *slog.Logger
	tracer    trace.Tracer
}

func NewOutboxRepository(txManager trm.TransactionManager, log *slog.Logger, tracer trace.Tracer) *OutboxRepository {
	return &OutboxRepository{
		txManager: txManager,
		logger:    log,
		tracer:    tracer,
	}
}

func (s *OutboxRepository) Create(ctx context.Context, msg outbox.Message) error {
	ctx, span := s.tracer.Start(ctx, "OutboxRepository.Create")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	args := db_queries.CreateOutboxMessageParams{
		Key:   msg.Key,
		Value: msg.Value,
		Topic: msg.Topic,
	}

	if err := querier.CreateOutboxMessage(ctx, args); err != nil {
		return err
	}

	return nil
}

func (s *OutboxRepository) CreateBatch(ctx context.Context, msgList []outbox.Message) error {
	ctx, span := s.tracer.Start(ctx, "OutboxRepository.CreateBatch")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	args := make([]db_queries.CreateOutboxMessagesParams, 0, len(msgList))

	for _, msg := range msgList {
		args = append(args, db_queries.CreateOutboxMessagesParams{
			Key:   msg.Key,
			Value: msg.Value,
			Topic: msg.Topic,
		})
	}

	if _, err := querier.CreateOutboxMessages(ctx, args); err != nil {
		return err
	}

	return nil
}

func (s *OutboxRepository) Read(ctx context.Context, topic string, limit int32) ([]outbox.Message, error) {
	ctx, span := s.tracer.Start(ctx, "OutboxRepository.Read")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	args := db_queries.ReadOutboxMessagesParams{
		Topic: topic,
		Limit: limit,
	}

	rows, err := querier.ReadOutboxMessages(ctx, args)
	if err != nil {
		return nil, err
	}

	msgList := make([]outbox.Message, 0, len(rows))

	for _, row := range rows {
		msg := outbox.Message{
			Key:   row.Key,
			Value: row.Value,
			Topic: row.Topic,
		}

		msgList = append(msgList, msg)
	}

	return msgList, nil
}
