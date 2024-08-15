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
	errChan := make(chan error)

	go func() {
		if err := s.srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := s.Shutdown(shutdownCtx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	listenAndServeErr := <-errChan
	shutdownErr := <-errChan

	if listenAndServeErr != nil {
		return listenAndServeErr
	}

	if shutdownErr != nil {
		return shutdownErr
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
