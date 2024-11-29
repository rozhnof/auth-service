package pgrepo

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/rozhnof/auth-service/internal/application/services"
	"github.com/rozhnof/auth-service/internal/domain/entities"
	db_queries "github.com/rozhnof/auth-service/internal/infrastructure/repository/queries"
	trm "github.com/rozhnof/auth-service/pkg/transaction_manager"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository struct {
	txManager trm.TransactionManager
	log       *slog.Logger
	tracer    trace.Tracer
}

func NewUserRepository(txManager trm.TransactionManager, log *slog.Logger, tracer trace.Tracer) *UserRepository {
	return &UserRepository{
		txManager: txManager,
		log:       log,
		tracer:    tracer,
	}
}

func (s *UserRepository) Create(ctx context.Context, user *entities.User) error {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Create")
	defer span.End()

	args := db_queries.CreateUserParams{
		ID:           user.ID(),
		Email:        user.Email(),
		HashPassword: user.Password().Hash(),
	}

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		if err := querier.CreateUser(ctx, args); err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return errors.Wrapf(services.ErrDuplicate, "user with email = %s already exists", user.Email())
				}
			}

			return err
		}

		if user.RefreshToken() != nil {
			refreshTokenArgs := db_queries.CreateOrUpdateRefreshTokenParams{
				UserID:       user.ID(),
				RefreshToken: user.RefreshToken().Token(),
				ExpiredAt:    user.RefreshToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRefreshToken(ctx, refreshTokenArgs); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserRepository) Update(ctx context.Context, user *entities.User) error {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Update")
	defer span.End()

	args := db_queries.UpdateUserParams{
		ID:           user.ID(),
		Email:        user.Email(),
		HashPassword: user.Password().Hash(),
	}

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		if user.RefreshToken() != nil {
			refreshTokenArgs := db_queries.CreateOrUpdateRefreshTokenParams{
				UserID:       user.ID(),
				RefreshToken: user.RefreshToken().Token(),
				ExpiredAt:    user.RefreshToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRefreshToken(ctx, refreshTokenArgs); err != nil {
				return err
			}
		} else {
			if err := querier.DeleteRefreshTokenByUserID(ctx, user.ID()); err != nil {
				return err
			}
		}

		if err := querier.UpdateUser(ctx, args); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(services.ErrObjectNotFound, "user with email = %s not exists", user.Email())
			}

			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "UserRepository.Delete")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		if err := querier.DeleteRefreshTokenByUserID(ctx, userID); err != nil {
			return err
		}

		if err := querier.DeleteUser(ctx, userID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRow, err := querier.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrapf(services.ErrObjectNotFound, "user with id = %s not exists", userID.String())
		}

		return nil, err
	}

	user := userRowToUser(userRow)

	return user, nil
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByEmail")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRow, err := querier.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrapf(services.ErrObjectNotFound, "user with email = %s not exists", email)
		}

		return nil, err
	}

	user := userRowToUser(userRow)

	return user, nil
}

func (s *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByRefreshToken")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRow, err := querier.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, services.ErrObjectNotFound
		}

		return nil, err
	}

	user := userRowToUser(userRow)

	return user, nil
}

func (s *UserRepository) List(ctx context.Context, filters *services.UserFilters, pagination *services.Pagination) ([]entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	var args db_queries.ListUserParams

	if filters != nil {
		args.UserIds = filters.UserIDs
	}

	if pagination != nil {
		args.Limit = &pagination.Limit
		args.Offset = pagination.Offset
	}

	db := s.txManager.TxOrDB(ctx)
	querier := db_queries.New(db)

	userRows, err := querier.ListUser(ctx, args)
	if err != nil {
		return nil, err
	}

	userList := make([]entities.User, 0, len(userRows))

	for _, userRow := range userRows {
		user := userRowToUser(userRow)

		userList = append(userList, *user)
	}

	return userList, nil
}
