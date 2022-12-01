package exporting

import (
	"context"

	"github.com/nikoksr/proji/pkg/packages/portability"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

type (
	// exporter is an interface for exporting a package to a specific format. The default implementation is TOML.
	exporter interface {
		ToConfig(ctx context.Context, pkg *domain.PackageExport, destination, fileType string) (string, error)
	}

	_exporter struct{}
)

var (
	// Compile-time check to ensure that _exporter implements exporter interface.
	_ exporter = (*_exporter)(nil)

	// std is the package-level default exporter.
	std = _exporter{}
)

// ToTOML creates a package config file at the specified destination. On success, the path to the file is returned.
func ToTOML(ctx context.Context, pkg *domain.PackageExport, destination string) (string, error) {
	return std.ToConfig(ctx, pkg, destination, portability.FileTypeTOML)
}

// ToJSON creates a package config file at the specified destination. On success, the path to the file is returned.
func ToJSON(ctx context.Context, pkg *domain.PackageExport, destination string) (string, error) {
	return std.ToConfig(ctx, pkg, destination, portability.FileTypeJSON)
}
