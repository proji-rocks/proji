package service

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

type projectService struct {
	timeout     time.Duration
	projectRepo domain.ProjectRepo
}

// Compile-time check to ensure that projectService implements the domain.ProjectService interface.
var _ domain.ProjectService = (*projectService)(nil)

// New returns a new instance of the project service. It requires a project repository.
func New(timeout time.Duration, repo domain.ProjectRepo) (domain.ProjectService, error) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	if repo == nil {
		return nil, errors.New("project repository is nil")
	}

	return &projectService{
		timeout:     timeout,
		projectRepo: repo,
	}, nil
}

// Fetch fetches a project from the repository.
func (p projectService) Fetch(ctx context.Context) ([]domain.Project, error) {
	return p.projectRepo.Fetch(ctx)
}

// GetByID fetches a project from the repository by id.
func (p projectService) GetByID(ctx context.Context, id string) (domain.Project, error) {
	return p.projectRepo.GetByID(ctx, id)
}

// Store stores a project in the repository.
func (p projectService) Store(ctx context.Context, _project *domain.ProjectAdd) error {
	return p.projectRepo.Store(ctx, _project)
}

// Update updates a project in the repository.
func (p projectService) Update(ctx context.Context, _project *domain.ProjectUpdate) error {
	return p.projectRepo.Update(ctx, _project)
}

// Remove removes a project from the repository.
func (p projectService) Remove(ctx context.Context, id string) error {
	return p.projectRepo.Remove(ctx, id)
}
