package health

import (
	"context"
	"net/http"

	"github.com/nikoksr/proji/pkg/sdk"
)

// Client is a client for the health API.
type Client struct {
	Backend *sdk.Backend
	Key     string
}

// NewClient creates a new Client for the health API (/api/v1/healthz).
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

// IsHealthy checks if the server is healthy.
func (c *Client) IsHealthy(ctx context.Context) (bool, error) {
	err := c.Backend.Call(ctx, http.MethodGet, "/api/v1/healthz", c.Key, nil, nil)

	return err == nil, err
}
