package platform

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/internal/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		host string
		err  error
	}{
		{
			name: "github.com",
			host: "github.com",
			err:  nil,
		},
		{
			name: "gitlab.com",
			host: "gitlab.com",
			err:  nil,
		},
		{
			name: "unknown.com",
			host: "unknown.com",
			err:  errors.New("unsupported platform \"unknown.com\""),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			platform, err := New(context.Background(), tc.host)
			if (err == nil) && platform == nil {
				t.Fatalf("expected platform, got nil")
			}
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected error %v, got %v", tc.err, err)
			}
		})
	}
}

func TestNewWithAuth(t *testing.T) {
	t.Parallel()

	cases := []struct {
		host string
		auth *config.Auth
		err  error
	}{
		{
			host: "github.com",
			auth: &config.Auth{
				GitHubToken: "token",
			},
			err: nil,
		},
		{
			host: "gitlab.com",
			auth: &config.Auth{
				GitLabToken: "token",
			},
			err: nil,
		},
		{
			host: "unknown.com",
			auth: &config.Auth{},
			err:  errors.New("unsupported platform \"unknown.com\""),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.host, func(t *testing.T) {
			t.Parallel()

			_, err := NewWithAuth(context.Background(), c.host, c.auth)
			if !errors.Is(err, c.err) {
				t.Fatalf("expected error %v, got %v", c.err, err)
			}
		})
	}
}
