package outbox

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/rozhnof/auth-service/internal/infrastructure/kafka"
	trm "github.com/rozhnof/auth-service/pkg/transaction_manager"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	txManager trm.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewRepository(txManager trm.TransactionManager, log *slog.Logger, tracer trace.Tracer) *Repository {
	return &Repository{
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *Repository) Create(ctx context.Context, message kafka.Message) error {
	ctx, span := s.tracer.Start(ctx, "Repository.Create")
	defer span.End()

	const query = `
		INSERT INTO outbox (
			key,
			value,
			topic
		) VALUES (
			$1, $2, $3
		);
	`

	db := s.txManager.TxOrDB(ctx)

	if _, err := db.Exec(ctx, query, message.Key, message.Value, message.Topic); err != nil {
		return err
	}

	return nil
}

func (s *Repository) CreateList(ctx context.Context, messages []kafka.Message) error {
	ctx, span := s.tracer.Start(ctx, "Repository.CreateList")
	defer span.End()

	rows := [][]interface{}{}
	for _, message := range messages {
		rows = append(rows, []interface{}{
			message.Key,
			message.Value,
			message.Topic,
		})
	}

	db := s.txManager.TxOrDB(ctx)

	if _, err := db.CopyFrom(ctx, pgx.Identifier{"outbox"}, []string{"key", "value", "topic"}, pgx.CopyFromRows(rows)); err != nil {
		return err
	}

	return nil
}

func (s *Repository) Read(ctx context.Context, topic string, limit int32) ([]kafka.Message, error) {
	ctx, span := s.tracer.Start(ctx, "Repository.Read")
	defer span.End()

	const query = `
		WITH messages AS (
			SELECT 
				id, 
				key, 
				value, 
				topic
			FROM 
				outbox o
			WHERE 
				o.topic = $1
				AND o.is_read = FALSE
				AND o.deleted_at IS NULL
			ORDER BY o.created_at
			LIMIT $2
		)
		UPDATE 
			outbox
		SET 
			is_read = TRUE
		FROM 
			messages m
		WHERE 
			outbox.id = m.id
		RETURNING 
			m.key, 
			m.value, 
			m.topic;
	`

	db := s.txManager.TxOrDB(ctx)

	rows, err := db.Query(ctx, query, topic, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	messageList, err := pgx.CollectRows(rows, pgx.RowToStructByName[kafka.Message])
	if err != nil {
		return nil, err
	}

	return messageList, nil
}
