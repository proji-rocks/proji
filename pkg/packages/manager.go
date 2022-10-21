package packages

import (
	"context"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

// Manager is an interface for managing packages.
type Manager interface {
	Fetch(ctx context.Context) ([]domain.Package, error)
	GetByLabel(ctx context.Context, label string) (domain.Package, error)
	Store(ctx context.Context, _package *domain.PackageAdd) error
	Update(ctx context.Context, _package *domain.PackageUpdate) error
	Remove(ctx context.Context, id string) error
	String() string
}
