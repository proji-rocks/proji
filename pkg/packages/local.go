package packages

import (
	"context"
	"path"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

type paths struct {
	base      string
	plugins   string
	templates string
}

// localManager is a local package manager. It manages packages in a local directory. It is used by the standalone proji
// binary. If you want to use the proji API, use the remoteManager instead.
type localManager struct {
	auth           *config.Auth
	paths          paths
	packageService domain.PackageService
}

// Compile-time check to ensure that localManager implements the Manager interface.
var _ Manager = &localManager{}

// SetBaseDirectory sets the base directory for the local package manager. Relative to the base directory are the
// plugins and templates directories. The default base directory is "proji".
func (m *localManager) SetBaseDirectory(dir string) {
	if m == nil || dir == "" {
		return
	}

	m.paths.base = dir
	m.paths.plugins = path.Join(dir, "plugins")
	m.paths.templates = path.Join(dir, "templates")
}

// NewLocalManager creates a new local package manager. It requires a domain.PackageService to be set. If you want to
// use the proji API, use the remoteManager instead. The localManager is used by the standalone proji binary and manages
// packages in a local directory. By default, the localManager uses the proji directory as base directory.
func NewLocalManager(auth *config.Auth, service domain.PackageService) (Manager, error) {
	if service == nil {
		return nil, errors.New("service is required")
	}
	if auth == nil {
		auth = &config.Auth{}
	}

	manager := &localManager{
		auth:           auth,
		packageService: service,
	}

	manager.SetBaseDirectory("proji")

	return manager, nil
}

// Fetch fetches all packages from the local storage.
func (m *localManager) Fetch(ctx context.Context) ([]domain.Package, error) {
	return m.packageService.Fetch(ctx)
}

// GetByLabel fetches a package from the local storage by its label.
func (m *localManager) GetByLabel(ctx context.Context, label string) (domain.Package, error) {
	return m.packageService.GetByLabel(ctx, label)
}

// Store stores a package on the local storage.
func (m *localManager) Store(ctx context.Context, _package *domain.PackageAdd) error {
	if err := m.downloadDependencies(ctx, _package); err != nil {
		return errors.Wrap(err, "download dependencies")
	}

	return m.packageService.Store(ctx, _package)
}

// Update updates a package on the local storage.
func (m *localManager) Update(ctx context.Context, _package *domain.PackageUpdate) error {
	return m.packageService.Update(ctx, _package)
}

// Remove removes a package from the local storage.
func (m *localManager) Remove(ctx context.Context, id string) error {
	return m.packageService.Remove(ctx, id)
}

// String returns the name of the local package manager - "local".
func (m *localManager) String() string {
	return "local"
}
