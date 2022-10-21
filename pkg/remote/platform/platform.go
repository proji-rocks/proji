package platform

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/nikoksr/proji/pkg/remote/platform/github"
	"github.com/nikoksr/proji/pkg/remote/platform/gitlab"
)

// NewWithAuth returns a platform for the given URL. Unsupported platforms will result in an error. The authentication
// is required for private repositories.
//
// The supported platforms are:
//   - github.com
//   - gitlab.com
func NewWithAuth(ctx context.Context, host string, auth *config.Auth) (platform remote.Platform, err error) {
	if auth == nil {
		auth = &config.Auth{}
	}

	switch host {
	case remote.HostGitHub:
		platform = github.New(ctx, auth.GitHubToken)
	case remote.HostGitLab:
		platform, err = gitlab.New(ctx, auth.GitLabToken)
	default:
		err = errors.Errorf("unsupported platform %q", host)
	}

	if err != nil {
		return nil, errors.Wrap(err, "create platform")
	}

	return platform, nil
}

// New returns a platform for the given URL. Unsupported platforms will result in an error. It uses no authentication,
// so it is not possible to access private repositories and rate-limit might be easily reached.
//
// The supported platforms are:
//   - github.com
//   - gitlab.com
func New(ctx context.Context, host string) (platform remote.Platform, err error) {
	return NewWithAuth(ctx, host, nil)
}
