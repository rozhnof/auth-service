package auth

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_service_requests_total",
			Help: "Total number of requests processed by the Auth Service.",
		},
		[]string{"path", "method", "status"},
	)

	errorRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_service_errors_requests_total",
			Help: "Total number of error requests processed by the Auth Service.",
		},
		[]string{"path", "method", "status"},
	)
)
