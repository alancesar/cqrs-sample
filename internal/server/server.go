package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type (
	Server struct {
		handler http.Handler
	}
)

func New(handler http.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) StartWithGracefulShutdown(ctx context.Context, addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: s.handler,
	}

	errs := make(chan error)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()

	log.Println("started", addr)
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errs:
			return err
		case <-notifyCtx.Done():
			log.Println("shutting down", addr)
			stop()

			withTimeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			if err := server.Shutdown(withTimeoutCtx); err != nil {
				cancel()
				return err
			}

			cancel()
			log.Println("stopped", addr)
			return nil
		}
	}
}
