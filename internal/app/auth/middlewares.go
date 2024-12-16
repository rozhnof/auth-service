package auth

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		if c.Writer.Status() >= 500 {
			log.InfoContext(
				c.Request.Context(),
				"internal server error",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Int("status", c.Writer.Status()),
				slog.String("address", c.Request.RemoteAddr),
				slog.String("duration", duration.String()),
			)
		} else if duration > time.Second {
			log.InfoContext(
				c.Request.Context(),
				"long time response",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Int("status", c.Writer.Status()),
				slog.String("address", c.Request.RemoteAddr),
				slog.String("duration", duration.String()),
			)
		}
	}
}

func PrometheusMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(serviceName, path))

		c.Next()

		timer.ObserveDuration()

		status := c.Writer.Status()

		responseStatus.WithLabelValues(serviceName, strconv.Itoa(status)).Inc()
		totalRequests.WithLabelValues(serviceName, path).Inc()
	}
}
