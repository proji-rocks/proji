package packages

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/sdk"
)

// Client is a client for the packages API.
type Client struct {
	Backend *sdk.Backend
	Key     string
}

// NewClient creates a new Client for the packages API (/api/v1/packages).
func NewClient(serverURL string) (*Client, error) {
	client, err := sdk.NewBackend(serverURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Backend: client,
		Key:     "",
	}, nil
}

// GetByLabel fetches a package from the remote server by its label.
func (c *Client) GetByLabel(ctx context.Context, label string) (domain.Package, error) {
	if label == "" {
		return domain.Package{}, errors.New("label is required")
	}

	var pkg domain.Package
	err := c.Backend.Call(ctx, http.MethodGet, "/api/v1/packages/"+label, c.Key, nil, &pkg)

	return pkg, err
}

// Fetch fetches all packages from the remote server.
func (c *Client) Fetch(ctx context.Context) ([]domain.Package, error) {
	var packageList []domain.Package
	err := c.Backend.Call(ctx, http.MethodGet, "/api/v1/packages", c.Key, nil, &packageList)

	return packageList, err
}

// Store stores a package on the remote server.
func (c *Client) Store(ctx context.Context, pkg *domain.PackageAdd) error {
	return c.Backend.Call(ctx, http.MethodPost, "/api/v1/packages", c.Key, pkg, nil)
}

// Update updates a package on the remote server.
func (c *Client) Update(ctx context.Context, pkg *domain.PackageUpdate) error {
	return c.Backend.Call(ctx, http.MethodPut, "/api/v1/packages/"+pkg.Label, c.Key, pkg, nil)
}

// Remove deletes a package from the remote server.
func (c *Client) Remove(ctx context.Context, label string) error {
	return c.Backend.Call(ctx, http.MethodDelete, "/api/v1/packages/"+label, c.Key, nil, nil)
}
