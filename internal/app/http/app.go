package http_app

import (
	"auth/internal/auth/application/services"
	redis_cache "auth/internal/auth/infrastructure/cache/redis"
	postgres_session_repository "auth/internal/auth/infrastructure/repository/postgres/session"
	postgres_user_repository "auth/internal/auth/infrastructure/repository/postgres/user"
	http_handlers "auth/internal/auth/presentation/handlers/http"
	"auth/internal/pkg/config"
	postgres_database "auth/internal/pkg/database/postgres"
	"auth/internal/pkg/password_manager"
	http_server "auth/internal/pkg/server/http"
	"auth/internal/pkg/token_manager"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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

var tracer trace.Tracer

const (
	serviceName = "auth-service"
)

type App struct {
	server *http_server.HTTPServer
	log    *slog.Logger
	router *gin.Engine
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	postgresDatabase, err := InitPostgresDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	redisCache, err := InitRedisCache(ctx, cfg)
	if err != nil {
		return nil, err
	}

	shutdown, err := InitTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer shutdown(ctx)

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
		server: server,
		log:    log,
		router: router,
	}, nil
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

// func InitTracer(ctx context.Context, name string) (func(ctx context.Context) error, error) {
// 	exporter, err := otlptracehttp.New(
// 		ctx,
// 		otlptracehttp.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resources, err := resource.Merge(
// 		resource.Default(),
// 		resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceNameKey.String(serviceName),
// 		),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	provider := trace.NewTracerProvider(
// 		trace.WithBatcher(exporter),
// 		trace.WithResource(resources),
// 	)

// 	otel.SetTracerProvider(provider)

// 	otel.SetTextMapPropagator(propagation.TraceContext{})

// 	if err := exporter.Start(ctx); err != nil {
// 		return nil, errors.Wrap(err, "exporter error")
// 	}

// 	return provider.Shutdown, nil
// }

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func InitTracer(ctx context.Context) (shutdown func(context.Context) error, err error) {
	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer("example.io/package/name")

	go func() {
		for {
			time.Sleep(time.Second)
			tp.ForceFlush(ctx)
		}
	}()

	return tp.Shutdown, nil
}

func newExporter(ctx context.Context) (*stdouttrace.Exporter, error) {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
