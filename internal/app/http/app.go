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
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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

type App struct {
	server *http_server.HTTPServer
	log    *slog.Logger
	router *gin.Engine
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	httpServerConfig := http_server.Config{
		Address:         cfg.Server.HTTP.Address,
		ShutdownTimeout: cfg.Server.HTTP.ShutdownTimeout,
		TLSConfig:       nil,
	}

	databaseConfig := postgres_database.Config{
		Address:  cfg.Repository.Address,
		Port:     cfg.Repository.Port,
		User:     cfg.Repository.User,
		Password: cfg.Repository.Password,
		DB:       cfg.Repository.DB,
		SSL:      cfg.Repository.SSL,
	}

	postgresDatabase, err := postgres_database.NewDatabase(ctx, databaseConfig)
	if err != nil {
		return nil, err
	}

	redisConfig := redis_cache.RedisConfig{
		Address:      cfg.Cache.Redis.Address,
		Port:         cfg.Cache.Redis.Port,
		User:         cfg.Cache.Redis.User,
		Password:     cfg.Cache.Redis.Password,
		UserPassword: cfg.Cache.Redis.UserPassword,
		DB:           cfg.Cache.Redis.DB,
	}

	redis, err := redis_cache.NewRedis(ctx, redisConfig)
	if err != nil {
		return nil, err
	}

	var (
		userCache    = redis_cache.NewUserCache(redis)
		sessionCache = redis_cache.NewSessionCache(redis)
	)

	var (
		transactionManager = postgres_database.NewTransactionManager(postgresDatabase)
		userRepository     = postgres_user_repository.NewUserRepository(transactionManager, log)
		sessionRepository  = postgres_session_repository.NewSessionRepository(transactionManager, log)
		atManager          = token_manager.NewAccessTokenManager(cfg.Service.Tokens.Access.Timeout, []byte(cfg.Service.Tokens.Access.SecretKey))
		rtManager          = token_manager.NewRefreshTokenManager(cfg.Service.Tokens.Refresh.Timeout)
		passwordManager    = password_manager.NewPasswordManager()
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

	userService, err := services.NewUserService(userServiceDependencies, log)
	if err != nil {
		return nil, err
	}

	authHandler := http_handlers.NewAuthHandler(userService, log)
	router := gin.New()

	InitRoutes(router.Use(PrometheusMiddleware()), authHandler)

	// Init monitoring
	//prometheus.MustRegister(requestsTotal, requestDuration)
	//go http.ListenAndServe(":9091", promhttp.Handler())

	server := http_server.New(httpServerConfig, router)

	return &App{
		server: server,
		log:    log,
		router: router,
	}, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		method := c.Request.Method
		elapsed := time.Since(start).Milliseconds()
		requestsTotal.WithLabelValues(method).Inc()
		requestDuration.WithLabelValues(method).Observe(float64(elapsed))
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting http server")

	return a.server.Run(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("stopping http server")

	return a.server.Shutdown(ctx)
}
