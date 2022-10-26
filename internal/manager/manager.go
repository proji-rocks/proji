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

// NewPackageManager is a convenience function that connects to a package manager based on the given address. If the
// address is empty, it will connect to the local package manager. Otherwise, it will connect to the remote package
// manager.
func NewPackageManager(ctx context.Context, address string, db *database.DB, auth *config.Auth) (packages.Manager, error) {
	logger := simplog.FromContext(ctx)

	// If an address is given, interpret that as an intent to connect to a remote package manager.
	logger.Debugf("creating a package manager")

	address = strings.TrimSpace(address)
	if address != "" {
		logger.Debugf("server address not empty, connecting to remote package manager at %s", address)

		return packages.NewRemoteManager(address)
	}

	// Otherwise, connect to the local package manager.
	logger.Debugf("server address is empty, creating a local package manager")

	repo, err := packageRepo.New(db)
	if err != nil {
		return nil, errors.Wrap(err, "create package repository")
	}

	service, err := packageService.New(defaultServiceTimeout, repo)
	if err != nil {
		return nil, errors.Wrap(err, "create package service")
	}

	// Create the local package manager.
	return packages.NewLocalManager(auth, service)
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
