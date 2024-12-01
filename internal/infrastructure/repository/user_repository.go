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

		if user.RefreshToken() != nil && user.RefreshToken().Valid() {
			refreshTokenArgs := db_queries.CreateOrUpdateRefreshTokenParams{
				UserID:    user.ID(),
				Token:     user.RefreshToken().Token(),
				ExpiredAt: user.RefreshToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRefreshToken(ctx, refreshTokenArgs); err != nil {
				return err
			}
		}

		if user.RegisterToken() != nil && user.RegisterToken().Valid() {
			registerTokenArgs := db_queries.CreateOrUpdateRegisterTokenParams{
				UserID:    user.ID(),
				Token:     user.RegisterToken().Token(),
				ExpiredAt: user.RegisterToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRegisterToken(ctx, registerTokenArgs); err != nil {
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

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		if user.RefreshToken() != nil && user.RefreshToken().Valid() {
			refreshTokenArgs := db_queries.CreateOrUpdateRefreshTokenParams{
				UserID:    user.ID(),
				Token:     user.RefreshToken().Token(),
				ExpiredAt: user.RefreshToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRefreshToken(ctx, refreshTokenArgs); err != nil {
				return err
			}
		} else {
			if err := querier.DeleteRefreshTokenByUserID(ctx, user.ID()); err != nil {
				return err
			}
		}

		if user.RegisterToken() != nil && user.RegisterToken().Valid() {
			registerTokenArgs := db_queries.CreateOrUpdateRegisterTokenParams{
				UserID:    user.ID(),
				Token:     user.RegisterToken().Token(),
				ExpiredAt: user.RegisterToken().ExpiredAt(),
			}

			if err := querier.CreateOrUpdateRegisterToken(ctx, registerTokenArgs); err != nil {
				return err
			}
		} else {
			if err := querier.DeleteRegisterTokenByUserID(ctx, user.ID()); err != nil {
				return err
			}
		}

		userArgs := db_queries.UpdateUserParams{
			ID:           user.ID(),
			Email:        user.Email(),
			HashPassword: user.Password().Hash(),
			Confirmed:    user.Confirmed(),
		}

		if err := querier.UpdateUser(ctx, userArgs); err != nil {
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

		if err := querier.DeleteRegisterTokenByUserID(ctx, userID); err != nil {
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

	var (
		userDTO          db_queries.User
		refreshTokenDTO  *db_queries.RefreshToken
		registerTokenDTO *db_queries.RegisterToken
	)

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		userRow, err := querier.GetUserByID(ctx, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(services.ErrObjectNotFound, "user with id = %s not exists", userID.String())
			}

			return err
		}

		userDTO = userRow.User

		refreshTokenRow, err := querier.GetRefreshTokenByUserID(ctx, userID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			refreshTokenDTO = &refreshTokenRow.RefreshToken
		}

		registerTokenRow, err := querier.GetRegisterTokenByUserID(ctx, userID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			registerTokenDTO = &registerTokenRow.RegisterToken
		}

		return nil
	}); err != nil {
		return nil, err
	}

	user := dtoToUser(userDTO, refreshTokenDTO, registerTokenDTO)

	return user, nil
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByEmail")
	defer span.End()

	var (
		userDTO          db_queries.User
		refreshTokenDTO  *db_queries.RefreshToken
		registerTokenDTO *db_queries.RegisterToken
	)

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		userRow, err := querier.GetUserByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(services.ErrObjectNotFound, "user with email = %s not exists", email)
			}

			return err
		}

		userDTO = userRow.User

		refreshTokenRow, err := querier.GetRefreshTokenByUserID(ctx, userDTO.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			refreshTokenDTO = &refreshTokenRow.RefreshToken
		}

		registerTokenRow, err := querier.GetRegisterTokenByUserID(ctx, userDTO.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			registerTokenDTO = &registerTokenRow.RegisterToken
		}

		return nil
	}); err != nil {
		return nil, err
	}

	user := dtoToUser(userDTO, refreshTokenDTO, registerTokenDTO)

	return user, nil
}

func (s *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.GetByRefreshToken")
	defer span.End()

	var (
		userDTO          db_queries.User
		refreshTokenDTO  *db_queries.RefreshToken
		registerTokenDTO *db_queries.RegisterToken
	)

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		userRow, err := querier.GetUserByRefreshToken(ctx, refreshToken)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(services.ErrObjectNotFound, "user with refresh token = %s not exists", refreshToken)
			}

			return err
		}

		userDTO = userRow.User

		refreshTokenRow, err := querier.GetRefreshTokenByUserID(ctx, userDTO.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			refreshTokenDTO = &refreshTokenRow.RefreshToken
		}

		registerTokenRow, err := querier.GetRegisterTokenByUserID(ctx, userDTO.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			registerTokenDTO = &registerTokenRow.RegisterToken
		}

		return nil
	}); err != nil {
		return nil, err
	}

	user := dtoToUser(userDTO, refreshTokenDTO, registerTokenDTO)

	return user, nil
}

func (s *UserRepository) List(ctx context.Context, filters *services.UserFilters, pagination *services.Pagination) ([]entities.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserRepository.List")
	defer span.End()

	var (
		userRows          []db_queries.ListUserRow
		refreshTokenRows  []db_queries.ListRefreshTokenRow
		registerTokenRows []db_queries.ListRegisterTokenRow
	)

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := db_queries.New(db)

		var userArgs db_queries.ListUserParams

		if filters != nil {
			userArgs.UserIds = filters.UserIDs
		}

		if pagination != nil {
			userArgs.Limit = &pagination.Limit
			userArgs.Offset = pagination.Offset
		}

		var err error

		userRows, err = querier.ListUser(ctx, userArgs)
		if err != nil {
			return err
		}

		refreshTokenArgs := db_queries.ListRefreshTokenParams{
			UserIds: userArgs.UserIds,
		}

		refreshTokenRows, err = querier.ListRefreshToken(ctx, refreshTokenArgs)
		if err != nil {
			return err
		}

		registerTokenArgs := db_queries.ListRegisterTokenParams{
			UserIds: userArgs.UserIds,
		}

		registerTokenRows, err = querier.ListRegisterToken(ctx, registerTokenArgs)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	userMap := make(
		map[uuid.UUID]struct {
			user          db_queries.User
			refreshToken  *db_queries.RefreshToken
			registerToken *db_queries.RegisterToken
		},
		len(userRows),
	)

	for _, userRow := range userRows {
		userID := userRow.User.ID

		userDTO := userMap[userID]
		userDTO.user = userRow.User

		userMap[userID] = userDTO
	}

	for _, refreshTokenRow := range refreshTokenRows {
		userID := refreshTokenRow.RefreshToken.UserID

		userDTO := userMap[userID]
		userDTO.refreshToken = &refreshTokenRow.RefreshToken

		userMap[userID] = userDTO
	}

	for _, registerTokenRow := range registerTokenRows {
		userID := registerTokenRow.RegisterToken.UserID

		userDTO := userMap[userID]
		userDTO.registerToken = &registerTokenRow.RegisterToken

		userMap[userID] = userDTO
	}

	userList := make([]entities.User, 0, len(userMap))
	for _, userDTO := range userMap {
		user := dtoToUser(userDTO.user, userDTO.refreshToken, userDTO.registerToken)
		userList = append(userList, *user)
	}

	return userList, nil
}
