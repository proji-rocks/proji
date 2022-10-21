package importing

import (
	"context"
	"regexp"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

type (
	// importer is the interface for importing packages. We use it to create packages from all kinds of different sources.
	importer interface {
		LocalPackage(ctx context.Context, path string) (*domain.PackageAdd, error)
		RemotePackage(ctx context.Context, url string) (*domain.PackageAdd, error)
		LocalFolderAsPackage(ctx context.Context, path string, exclude *regexp.Regexp) (*domain.PackageAdd, error)
		RepositoryAsPackage(ctx context.Context, url string, exclude *regexp.Regexp) (*domain.PackageAdd, error)
	}

	_importer struct{}
)

var (
	// Compile-time check to ensure that _importer implements importer interface.
	_ importer = (*_importer)(nil)

	// std is the package-level default importer.
	std = _importer{}
)

// LocalPackage creates a package from a local config file.
func LocalPackage(ctx context.Context, path string) (*domain.PackageAdd, error) {
	return std.LocalPackage(ctx, path)
}

// RemotePackage creates a package from a remote config file.
func RemotePackage(ctx context.Context, url string) (*domain.PackageAdd, error) {
	return std.RemotePackage(ctx, url)
}

// LocalFolderAsPackage creates a package from a local folder. It mimics the folders structure.
func LocalFolderAsPackage(ctx context.Context, path string, exclude *regexp.Regexp) (*domain.PackageAdd, error) {
	return std.LocalFolderAsPackage(ctx, path, exclude)
}

// RepositoryAsPackage creates a package from a repository. It mimics the repositories structure.
func RepositoryAsPackage(ctx context.Context, url string, exclude *regexp.Regexp) (*domain.PackageAdd, error) {
	return std.RepositoryAsPackage(ctx, url, exclude)
}
