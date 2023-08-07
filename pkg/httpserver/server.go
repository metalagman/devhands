package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

type ServerOption func(*Server)

func WithListenAddr(a string) ServerOption {
	return func(s *Server) {
		s.server.Addr = a
	}
}

func WithHandler(h http.Handler) ServerOption {
	return func(s *Server) {
		s.server.Handler = h
	}
}

func WithReadTimeout(t time.Duration) ServerOption {
	return func(s *Server) {
		s.server.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) ServerOption {
	return func(s *Server) {
		s.server.WriteTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) ServerOption {
	return func(s *Server) {
		s.server.IdleTimeout = t
	}
}

func New(opts ...ServerOption) *Server {
	s := &Server{
		server: &http.Server{},
	}

	for _, o := range opts {
		o(s)
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return s
}

func (s *Server) Stop() error {
	ctxShutdown, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	if err := s.server.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
