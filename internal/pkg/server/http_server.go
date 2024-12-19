package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
)

type HTTPServer struct {
	srv    http.Server
	logger *slog.Logger
}

func NewHTTPServer(ctx context.Context, address string, handler http.Handler, logger *slog.Logger) *HTTPServer {
	s := &HTTPServer{
		srv: http.Server{
			Addr:        address,
			Handler:     handler,
			BaseContext: func(net.Listener) context.Context { return ctx },
		},
		logger: logger,
	}

	return s
}

func (s *HTTPServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		if err := s.srv.Shutdown(ctx); err != nil {
			s.logger.Warn("failed shutdown http server", slog.String("error", err.Error()))
		}
	}()

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
