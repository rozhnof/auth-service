package auth

import (
	"log/slog"
	"net/http"
	"strconv"
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
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()

		requestCount.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
		if status >= http.StatusInternalServerError {
			errorRequestCount.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
		}
	}
}
