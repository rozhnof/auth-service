package http_server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

type Config struct {
	Address         string
	ShutdownTimeout time.Duration
	TLSConfig       *tls.Config
}

type HTTPServer struct {
	srv *http.Server
	cfg Config
}

func New(cfg Config, handler http.Handler) *HTTPServer {
	s := &HTTPServer{
		srv: &http.Server{
			Addr:      cfg.Address,
			TLSConfig: cfg.TLSConfig,
			Handler:   handler,
		},
		cfg: cfg,
	}

	return s
}

func (s *HTTPServer) Run(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) RunTLS(ctx context.Context) error {
	return s.srv.ListenAndServeTLS("", "")
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
