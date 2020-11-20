package projectservice

import (
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/nikoksr/proji/pkg/template_engine"
)

type projectService struct {
	projectStore   domain.ProjectStore
	templateEngine *template_engine.TemplateEngine
}

func New(store domain.ProjectStore, templateStartTag, templateEndTag string) domain.ProjectService {
	return &projectService{
		projectStore:   store,
		templateEngine: template_engine.NewTemplateEngine(templateStartTag, templateEndTag),
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
