package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rozhnof/auth-service/internal/application/services"
	"github.com/rozhnof/auth-service/internal/infrastructure/database/postgres"
	"github.com/rozhnof/auth-service/internal/infrastructure/database/redis"
	"github.com/rozhnof/auth-service/internal/infrastructure/kafka"
	pgrepo "github.com/rozhnof/auth-service/internal/infrastructure/repository"
	"github.com/rozhnof/auth-service/internal/infrastructure/secrets"
	"github.com/rozhnof/auth-service/internal/pkg/config"
	"github.com/rozhnof/auth-service/internal/pkg/outbox"
	"github.com/rozhnof/auth-service/internal/pkg/server"
	"github.com/rozhnof/auth-service/internal/presentation/clients"
	"github.com/rozhnof/auth-service/internal/presentation/handlers"
	trm "github.com/rozhnof/auth-service/pkg/transaction_manager"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ServiceName     = "Auth Service"
	OutboxBatchSize = 100
	OutboxInterval  = time.Second * 30
)

const (
	loginsTopic    = "logins"
	registersTopic = "registers"
)

const (
	googleCallbackURL = "http://localhost:8080/auth/google/callback"
)

var (
	googleOAuthScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

type Config struct {
	Mode     string         `yaml:"mode"    env-required:"true"`
	Server   config.Server  `yaml:"server"  env-required:"true"`
	Logger   config.Logger  `yaml:"logging" env-required:"true"`
	Tokens   config.Tokens  `yaml:"tokens"  env-required:"true"`
	Tracing  config.Tracing `yaml:"tracing" env-required:"true"`
	Kafka    config.Kafka   `yaml:"kafka"   env-required:"true"`
	Postgres config.Postgres
	Redis    config.Redis
}

type App struct {
	logger     *slog.Logger
	httpServer *server.HTTPServer
}

func NewApp(
	ctx context.Context,
	cfg *Config,
	logger *slog.Logger,
	postgresDatabase postgres.Database,
	redisDatabase redis.Database,
	kafkaProducer kafka.Producer,
	tracer trace.Tracer,
) *App {
	var (
		txManager     = trm.NewTransactionManager(postgresDatabase.Pool)
		secretManager = secrets.NewEnvSecretManager()
	)

	var (
		userRepository   = pgrepo.NewUserRepository(txManager, logger, tracer)
		outboxRepository = pgrepo.NewOutboxRepository(txManager, logger, tracer)
	)

	// var (
	// 	loginsKafkaSender    = kafka.NewMessageSender(kafkaProducer, loginsTopic)
	// 	registersKafkaSender = kafka.NewMessageSender(kafkaProducer, registersTopic)
	// )

	var (
		loginsOutboxSender    = outbox.NewMessageSender(outboxRepository, loginsTopic)
		registersOutboxSender = outbox.NewMessageSender(outboxRepository, registersTopic)
	)

	var (
		authServiceConfig = services.AuthServiceConfig{
			AccessTokenTTL:  cfg.Tokens.AccessTokenTTL,
			RefreshTokenTTL: cfg.Tokens.RefreshTokenTTL,
		}

		authService = services.NewAuthService(
			userRepository,
			txManager,
			secretManager,
			loginsOutboxSender,
			registersOutboxSender,
			logger,
			tracer,
			authServiceConfig,
		)
	)

	googleAuthHandlerConfig := oauth2.Config{
		RedirectURL:  googleCallbackURL,
		ClientID:     string(secretManager.GoogleClientID().Get()),
		ClientSecret: string(secretManager.GoogleClientSecret().Get()),
		Scopes:       googleOAuthScopes,
		Endpoint:     google.Endpoint,
	}

	googleAuthClient := clients.NewGoogleAuthClient(googleAuthHandlerConfig)

	var (
		authHandler       = handlers.NewAuthHandler(authService, logger, tracer)
		googleAuthHandler = handlers.NewGoogleAuthHandler(googleAuthClient, authService, logger, tracer)
	)

	gin.SetMode(cfg.Mode)

	var (
		router = gin.New()
	)

	router.Use(
		gin.Recovery(),
		otelgin.Middleware(ServiceName),
		LogMiddleware(logger),
		PrometheusMiddleware(ServiceName),
	)

	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(httpDuration)

	InitAuthRoutes(router, authHandler, googleAuthHandler)
	InitSwaggerRoutes(router)
	InitPrometheusRoutes(router)

	httpServer := server.NewHTTPServer(ctx, cfg.Server.Address, router, logger)

	return &App{
		logger:     logger,
		httpServer: httpServer,
	}
}

func (a *App) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}
