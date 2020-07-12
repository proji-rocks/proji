package projectservice

import (
	"github.com/nikoksr/proji/pkg/domain"
)

type projectService struct {
	projectStore domain.ProjectStore
}

func New(store domain.ProjectStore) domain.ProjectService {
	return &projectService{
		projectStore: store,
	}
}

func (ps projectService) StoreProject(p *domain.Project) error {
	return ps.projectStore.StoreProject(p)
}

func (ps projectService) LoadProject(path string) (*domain.Project, error) {
	return ps.projectStore.LoadProject(path)
}

func (ps projectService) LoadProjectList(paths ...string) ([]*domain.Project, error) {
	return ps.projectStore.LoadProjectList(paths...)
}

func (ps projectService) UpdateProjectLocation(oldPath, newPath string) error {
	return ps.UpdateProjectLocation(oldPath, newPath)
}

func (ps projectService) RemoveProject(path string) error {
	return ps.projectStore.RemoveProject(path)
}

func (ps projectService) PurgeProject(path string) error {
	return ps.projectStore.PurgeProject(path)
}
