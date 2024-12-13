package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()

		RequestCount.WithLabelValues(path, http.StatusText(status)).Inc()
		if status >= 400 {
			ErrorRequestCount.WithLabelValues(path, http.StatusText(status)).Inc()
		}
	}
}
