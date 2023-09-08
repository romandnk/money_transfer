package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
	notify          chan error
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:           defaultAddr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
	}

	srv := &Server{
		srv:             httpServer,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.srv.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
