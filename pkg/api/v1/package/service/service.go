package service

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/nikoksr/proji/pkg/remote/platform"
)

type packageService struct {
	timeout     time.Duration
	packageRepo domain.PackageRepo
}

// Compile-time check to ensure that packageService implements the domain.PackageService interface.
var _ domain.PackageService = (*packageService)(nil)

// New returns a new instance of the package service. It requires a package repository.
func New(timeout time.Duration, repo domain.PackageRepo) (domain.PackageService, error) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	if repo == nil {
		return nil, errors.New("package repository is nil")
	}

	return &packageService{
		timeout:     timeout,
		packageRepo: repo,
	}, nil
}

// Fetch fetches a package from the repository.
func (p packageService) Fetch(ctx context.Context) ([]domain.Package, error) {
	return p.packageRepo.Fetch(ctx)
}

// GetByLabel fetches a package from the repository by label.
func (p packageService) GetByLabel(ctx context.Context, label string) (domain.Package, error) {
	return p.packageRepo.GetByLabel(ctx, label)
}

// Store stores a package in the repository.
func (p packageService) Store(ctx context.Context, _package *domain.PackageAdd) error {
	return p.packageRepo.Store(ctx, _package)
}

// Update updates a package in the repository.
func (p packageService) Update(ctx context.Context, _package *domain.PackageUpdate) error {
	return p.packageRepo.Update(ctx, _package)
}

// UpdateFromUpstream updates a package from upstream. This usually means pulling the latest version from GitHub or
// GitLab.
func (p packageService) UpdateFromUpstream(ctx context.Context, _package *domain.PackageUpdate) error {
	if _package == nil {
		return errors.New("package is nil")
	}
	if _package.UpstreamURL == nil {
		return errors.New("package has no upstream url")
	}

	upstreamURL, err := remote.ParseRepoURL(*_package.UpstreamURL)
	if err != nil {
		return errors.Wrap(err, "parse upstream url")
	}

	// TODO: Uses no authentication.
	_, err = platform.New(ctx, upstreamURL.Hostname())
	if err != nil {
		return errors.Wrap(err, "identify platform")
	}

	return errors.New("not implemented")
}

// Remove removes a package from the repository.
func (p packageService) Remove(ctx context.Context, id string) error {
	return p.packageRepo.Remove(ctx, id)
}
