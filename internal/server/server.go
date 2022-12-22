package server

import (
	"context"
	"net/http"
	"time"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"moul.io/chizap"

	healthHandlers "github.com/nikoksr/proji/pkg/api/v1/health/delivery/http"
	packageHandlers "github.com/nikoksr/proji/pkg/api/v1/package/delivery/http"
	"github.com/nikoksr/proji/pkg/packages"
)

type (
	// Server is the main struct of the server. It contains the internal http server and the logger.
	Server struct {
		core               *http.Server
		logger             *zap.Logger
		isHealthy, isReady atomic.Bool
	}

	// Managers is a struct that contains all the managers that are used by the server. Currently only the package
	// manager is used. We're using this instead of passing the managers directly to the server because we want to be
	// able to easily add more managers in the future.
	Managers struct {
		Package packages.Manager
	}
)

const requestTimeout = 60 * time.Second

// newRouter creates a new chi router. It is required that the logger is already initialized.
func newRouter(logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.CleanPath)

	router.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(requestTimeout))

	return router
}

// newHTTPServer creates a new http server. It is required that the router is already initialized.
func newHTTPServer() *http.Server {
	return &http.Server{
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (s *Server) registerHandlers(managers *Managers) error {
	router := newRouter(s.logger)

	// Basic handlers
	healthHandlers.Register(s.logger, router, &s.isHealthy, &s.isReady)

	// Packages
	packageHandlers.Register(s.logger, router, managers.Package)

	// Set the router
	s.core.Handler = router

	return nil
}

// New creates a new Server. It uses sane defaults for the logger and router. The server can be started with the
// Run method.
func New(ctx context.Context, managers *Managers) (*Server, error) {
	logger := simplog.FromContext(ctx)

	if managers == nil {
		return nil, errors.New("no managers given to server")
	}

	logger.Info("setting up server")
	server := &Server{
		logger: logger.Desugar(),
		core:   newHTTPServer(),
	}

	// Register the handlers
	logger.Info("registering handlers")
	err := server.registerHandlers(managers)
	if err != nil {
		return nil, errors.Wrap(err, "register handlers")
	}

	return server, nil
}
