package http_app

import (
	"auth/internal/auth/application/services"
	redis_cache "auth/internal/auth/infrastructure/cache/redis"
	postgres_session_repository "auth/internal/auth/infrastructure/repository/postgres/session"
	postgres_user_repository "auth/internal/auth/infrastructure/repository/postgres/user"
	tracing "auth/internal/auth/infrastructure/tracer"
	http_handlers "auth/internal/auth/presentation/handlers/http"
	"auth/internal/pkg/config"
	postgres_database "auth/internal/pkg/database/postgres"
	"auth/internal/pkg/password_manager"
	http_server "auth/internal/pkg/server/http"
	"auth/internal/pkg/token_manager"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)
)

const (
	serviceName           = "Auth Service"
	tracerName            = "auth-service-tracer"
	traceExporterEndpoint = "localhost:4318"
)

type App struct {
	server  *http_server.HTTPServer
	log     *slog.Logger
	router  *gin.Engine
	closers []func(context.Context) error
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	closers := []func(context.Context) error{}

	postgresDatabase, err := InitPostgresDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}
	closers = append(closers, func(_ context.Context) error {
		postgresDatabase.Close()
		return nil
	})

	redisCache, err := InitRedisCache(ctx, cfg)
	if err != nil {
		return nil, err
	}
	closers = append(closers, redisCache.Close)

	tracer, shutdown, err := InitTracer(ctx, traceExporterEndpoint, tracerName, serviceName)
	if err != nil {
		return nil, err
	}
	closers = append(closers, shutdown)

	var (
		userCache    = redis_cache.NewUserCache(redisCache)
		sessionCache = redis_cache.NewSessionCache(redisCache)
	)

	var (
		transactionManager = postgres_database.NewTransactionManager(postgresDatabase)
	)

	var (
		userRepository    = postgres_user_repository.NewUserRepository(transactionManager, log)
		sessionRepository = postgres_session_repository.NewSessionRepository(transactionManager, log)
	)

	var (
		atManager       = token_manager.NewAccessTokenManager(cfg.Service.Tokens.Access.Timeout, []byte(cfg.Service.Tokens.Access.SecretKey))
		rtManager       = token_manager.NewRefreshTokenManager(cfg.Service.Tokens.Refresh.Timeout)
		passwordManager = password_manager.NewPasswordManager()
	)

	userServiceDependencies := services.Dependencies{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		TxManager:         transactionManager,
		AtManager:         atManager,
		RtManager:         rtManager,
		PasswordManager:   passwordManager,
		UserCache:         userCache,
		SessionCache:      sessionCache,
	}

	userService, err := services.NewUserService(userServiceDependencies, log, tracer)
	if err != nil {
		return nil, err
	}

	var (
		authHandler = http_handlers.NewAuthHandler(userService, log, tracer)
	)

	if cfg.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		return nil, fmt.Errorf("invalid mode: %s", cfg.Mode)
	}

	var (
		router = gin.New()
	)

	middlewareList := []gin.HandlerFunc{
		otelgin.Middleware(serviceName),
		PrometheusMiddleware(),
		LogMiddleware(log),
	}
	router.Use(middlewareList...)

	InitAuthRoutes(router, authHandler)

	// Init monitoring
	//prometheus.MustRegister(requestsTotal, requestDuration)
	//go http.ListenAndServe(":9091", promhttp.Handler())

	var (
		server = InitHTTPServer(ctx, cfg, router)
	)

	return &App{
		server:  server,
		log:     log,
		router:  router,
		closers: closers,
	}, nil
}

func (a *App) Close(ctx context.Context) error {
	var errs error

	for _, closer := range a.closers {
		if err := closer(ctx); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting http server")

	return a.server.Run(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("stopping http server")

	return a.server.Shutdown(ctx)
}

func InitPostgresDatabase(ctx context.Context, cfg *config.Config) (*postgres_database.Database, error) {
	databaseConfig := postgres_database.Config{
		Address:  cfg.Repository.Address,
		Port:     cfg.Repository.Port,
		User:     cfg.Repository.User,
		Password: cfg.Repository.Password,
		DB:       cfg.Repository.DB,
		SSL:      cfg.Repository.SSL,
	}

	return postgres_database.NewDatabase(ctx, databaseConfig)
}

func InitRedisCache(ctx context.Context, cfg *config.Config) (*redis_cache.Redis, error) {
	redisConfig := redis_cache.RedisConfig{
		Address:      cfg.Cache.Redis.Address,
		Port:         cfg.Cache.Redis.Port,
		User:         cfg.Cache.Redis.User,
		Password:     cfg.Cache.Redis.Password,
		UserPassword: cfg.Cache.Redis.UserPassword,
		DB:           cfg.Cache.Redis.DB,
	}

	return redis_cache.NewRedis(ctx, redisConfig)
}

func InitHTTPServer(ctx context.Context, cfg *config.Config, handler http.Handler) *http_server.HTTPServer {
	httpServerConfig := http_server.Config{
		Address:         cfg.Server.HTTP.Address,
		ShutdownTimeout: cfg.Server.HTTP.ShutdownTimeout,
		TLSConfig:       nil,
	}

	return http_server.New(ctx, httpServerConfig, handler)
}

func InitTracer(ctx context.Context, url string, tracerName string, serviceName string) (tracer trace.Tracer, shutdown func(context.Context) error, err error) {
	exporter, err := tracing.NewJaegerExporter(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	provider, err := tracing.NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, nil, err
	}

	otel.SetTracerProvider(provider)

	tracer = provider.Tracer(tracerName)

	return tracer, provider.Shutdown, nil
}
