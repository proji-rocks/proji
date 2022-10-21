package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *Server) serve(ctx context.Context) error {
	defer func() {
		s.isHealthy.Store(false)
		s.isReady.Store(false)
	}()

	// Prepare graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	// Start server in a separate goroutine and listen for context cancellation in
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.core.ListenAndServe()
	})
	group.Go(func() error {
		<-groupCtx.Done()
		return s.shutdown(context.Background())
	})

	// Set server as ready and healthy
	s.isHealthy.Store(true)
	s.isReady.Store(true)

	err := group.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Run the server. It blocks until the context is cancelled. It is required that the internal http server and router
// are already initialized. If the address is empty, it will listen on the default address (:8080).
// The server will be shutdown gracefully when the context is cancelled. It will automatically clean up after itself.
func (s *Server) Run(ctx context.Context, address string) error {
	// Sanity checks
	if s.core == nil {
		return errors.New("server not initialized")
	}
	if s.core.Handler == nil {
		return errors.New("router not initialized")
	}

	// Normalize address
	if address == "" {
		address = ":8080"
	}
	s.core.Addr = address

	s.logger.Info("starting server", zap.String("address", address))

	return s.serve(ctx)
}

// shutdown the server gracefully. It shuts down the internal http server and flushes the logs.
func (s *Server) shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")

	if s.core != nil {
		if err := s.core.Shutdown(ctx); err != nil {
			s.logger.Error("server.shutdown: ", zap.Error(err))
		}
	}

	err := s.logger.Sync()

	return errors.Wrap(err, "flush logger")
}
