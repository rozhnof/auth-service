package server

import (
	"context"
	"errors"
	"net"
	"net/http"
)

type HTTPServer struct {
	srv http.Server
}

func NewHTTPServer(ctx context.Context, address string, handler http.Handler) *HTTPServer {
	s := &HTTPServer{
		srv: http.Server{
			Addr:        address,
			Handler:     handler,
			BaseContext: func(net.Listener) context.Context { return ctx },
		},
	}

	return s
}

func (s *HTTPServer) Run() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *HTTPServer) Shutdown() error {
	return s.srv.Shutdown(context.Background())
}
