package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

type Config struct {
	Address         string
	ShutdownTimeout time.Duration
	TLSConfig       *tls.Config
}

type Server struct {
	srv *http.Server
	cfg Config
}

func New(cfg Config, handler http.Handler) *Server {
	s := &Server{
		srv: &http.Server{
			Addr:      cfg.Address,
			TLSConfig: cfg.TLSConfig,
			Handler:   handler,
		},
		cfg: cfg,
	}

	return s
}

func (s *Server) RunTLS(ctx context.Context) error {
	var listenAndServeErr error
	go func() {
		if err := s.srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			listenAndServeErr = err
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		if listenAndServeErr != nil {
			return errors.Join(err, listenAndServeErr)
		}
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
