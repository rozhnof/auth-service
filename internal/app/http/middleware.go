package http_app

import (
	"log/slog"
	"net/http"
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
		if c.Request.Method == http.MethodGet {
			c.Next()

			return
		}

		var (
			start = time.Now()
			path  = c.Request.URL.Path
			query = c.Request.URL.RawQuery
		)

		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		c.Next()

		var (
			end     = time.Now()
			status  = c.Writer.Status()
			method  = c.Request.Method
			host    = c.Request.Host
			route   = c.FullPath()
			latency = end.Sub(start)
		)

		requestAttributeList := []slog.Attr{
			slog.Time("time", start.UTC()),
			slog.String("method", method),
			slog.String("host", host),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.String("route", route),
		}

		responseAttributeList := []slog.Attr{
			slog.Time("time", end.UTC()),
			slog.Duration("latency", latency),
			slog.Int("status", status),
		}

		attributeList := []slog.Attr{
			{
				Key:   "request",
				Value: slog.GroupValue(requestAttributeList...),
			},
			{
				Key:   "response",
				Value: slog.GroupValue(responseAttributeList...),
			},
		}

		msg := "Incoming request"
		if status >= 500 {
			msg = c.Errors.String()
		} else if status >= 400 {
			msg = c.Errors.String()
		}

		log.LogAttrs(c.Request.Context(), slog.LevelInfo, msg, attributeList...)
	}
}
