package projects

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

// Manager is an interface for managing projects.
type Manager interface {
	Fetch(ctx context.Context) ([]domain.Project, error)
	GetByID(ctx context.Context, id string) (domain.Project, error)
	Store(ctx context.Context, project *domain.ProjectAdd) error
	Update(ctx context.Context, project *domain.ProjectUpdate) error
	Remove(ctx context.Context, id string) error
}

// manager is the default implementation of the Manager interface. In comparison to packages, projects will (at least
// for now) be stored in a local directory and never uploaded to a remote server. So the implementation details differ
// from the packages' implementation.
type manager struct {
	service domain.ProjectService
}

// Compile-time check to ensure that manager implements the Manager interface.
var _ Manager = &manager{}

// NewManager creates a new manager. It requires a domain.ProjectService to be set.
func NewManager(service domain.ProjectService) (Manager, error) {
	if service == nil {
		return nil, errors.New("service is required")
	}

	return &manager{
		service: service,
	}, nil
}

// Fetch fetches all projects from the local storage.
func (m *manager) Fetch(ctx context.Context) ([]domain.Project, error) {
	return m.service.Fetch(ctx)
}

// GetByID fetches a project by its ID.
func (m *manager) GetByID(ctx context.Context, id string) (domain.Project, error) {
	return m.service.GetByID(ctx, id)
}

// Store stores a project in the local storage.
func (m *manager) Store(ctx context.Context, project *domain.ProjectAdd) error {
	return m.service.Store(ctx, project)
}

// Update updates a project in the local storage.
func (m *manager) Update(ctx context.Context, project *domain.ProjectUpdate) error {
	return m.service.Update(ctx, project)
}

// Remove removes a project from the local storage.
func (m *manager) Remove(ctx context.Context, id string) error {
	return m.service.Remove(ctx, id)
}
