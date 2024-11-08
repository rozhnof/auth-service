package http_app

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

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

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		log = log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("address", c.Request.RemoteAddr),
			slog.String("duration", duration.String()),
		)

		if c.Writer.Status() >= 500 {
			log.Info("internal server error")
		} else if duration > time.Second {
			log.Info("long time response")
		}
	}
}
