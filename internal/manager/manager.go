package manager

import (
	"context"
	"strings"
	"time"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/internal/config"
	packageRepo "github.com/nikoksr/proji/pkg/api/v1/package/repository/bolt"
	packageService "github.com/nikoksr/proji/pkg/api/v1/package/service"
	projectRepo "github.com/nikoksr/proji/pkg/api/v1/project/repository/bolt"
	projectService "github.com/nikoksr/proji/pkg/api/v1/project/service"
	database "github.com/nikoksr/proji/pkg/database/bolt"
	"github.com/nikoksr/proji/pkg/packages"
	"github.com/nikoksr/proji/pkg/projects"
)

// TODO: Should probably be configurable.
const defaultServiceTimeout = 5 * time.Second

type LocalPaths struct {
	Base      string
	Plugins   string
	Templates string
}

// Config is a package manager configuration. This config is shared between different types of package managers.
type Config struct {
	// Address is the address of the remote package manager. If empty, the local package manager will be used.
	Address string

	// Auth contains the authentication information for the remote package manager.
	Auth *config.Auth

	// DB is the database connection.
	DB *database.DB

	// LocalPaths contains local filesystem paths that point to the base directory, the plugins/ directory and the
	// templates/ directory. These paths are used by the local package manager to persist packages and templates.
	LocalPaths *LocalPaths
}

// NewPackageManager is a convenience function that connects to a package manager based on the given address. If the
// address is empty, it will connect to the local package manager. Otherwise, it will connect to the remote package
// manager.
func NewPackageManager(ctx context.Context, config Config) (packages.Manager, error) {
	if config == (Config{}) {
		return nil, errors.New("config is required")
	}

	logger := simplog.FromContext(ctx)
	logger.Debugf("creating a package manager")

	// If an address is given, interpret that as an intent to connect to a remote package manager.
	config.Address = strings.TrimSpace(config.Address)
	if config.Address != "" {
		logger.Debugf("server address not empty, connecting to remote package manager at %q", config.Address)

		return packages.NewRemoteManager(config.Address)
	}

	// Otherwise, connect to the local package manager.
	logger.Debugf("server address is empty, creating a local package manager")

	repo, err := packageRepo.New(config.DB)
	if err != nil {
		return nil, errors.Wrap(err, "create package repository")
	}

	service, err := packageService.New(defaultServiceTimeout, repo)
	if err != nil {
		return nil, errors.Wrap(err, "create package service")
	}

	// Create the local package manager.
	return packages.NewLocalManager(config.Auth, service)
}

// NewProjectManager returns a new project manager. Compared to the package manager, the project manager is always local,
// at least for now.
func NewProjectManager(ctx context.Context, db *database.DB) (projects.Manager, error) {
	logger := simplog.FromContext(ctx)

	// Create the project manager.
	logger.Debugf("creating a project manager")

	repo, err := projectRepo.New(db)
	if err != nil {
		return nil, errors.Wrap(err, "create project repository")
	}

	service, err := projectService.New(defaultServiceTimeout, repo)
	if err != nil {
		return nil, errors.Wrap(err, "create project service")
	}

	// Create the local project manager.
	return projects.NewManager(service)
}
