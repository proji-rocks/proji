package packages

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/sdk/packages"
)

// remoteManager is a remote package manager. It manages packages in a remote proji server. The server itself most
// likely uses a localManager and offers an API to control it.
type remoteManager struct {
	client *packages.Client
}

// NewRemoteManager creates a new remote package manager. It requires a server URL. It manages packages in a remote
// proji server. The server itself most likely uses a localManager and offers an API to control it.
func NewRemoteManager(serverURL string) (Manager, error) {
	// Create packages API client.
	client, err := packages.NewClient(serverURL)
	if err != nil {
		return nil, errors.Wrap(err, "create API client")
	}

	// Create manager
	m := &remoteManager{
		client: client,
	}

	return m, nil
}

// Compile-time check to ensure that remoteManager implements the Manager interface.
var _ Manager = &remoteManager{}

// Fetch fetches all packages from the remote server.
func (m *remoteManager) Fetch(ctx context.Context) ([]domain.Package, error) {
	return m.client.Fetch(ctx)
}

// GetByLabel fetches a package from the remote server by its label.
func (m *remoteManager) GetByLabel(ctx context.Context, label string) (domain.Package, error) {
	return m.client.GetByLabel(ctx, label)
}

// Store stores a package on the remote server.
func (m *remoteManager) Store(ctx context.Context, _package *domain.PackageAdd) error {
	return m.client.Store(ctx, _package)
}

// Update updates a package on the remote server.
func (m *remoteManager) Update(ctx context.Context, pkg *domain.PackageUpdate) error {
	return m.client.Update(ctx, pkg)
}

// Remove removes a package from the remote server.
func (m *remoteManager) Remove(ctx context.Context, label string) error {
	return m.client.Remove(ctx, label)
}

// String returns the name of the remote package manager - "remote".
func (m *remoteManager) String() string {
	return "remote"
}
