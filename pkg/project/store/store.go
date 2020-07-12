package projectstore

import (
	"github.com/nikoksr/proji/pkg/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.ProjectStore {
	return &projectStore{
		db: db,
	}
}

func (ps *projectStore) StoreProject(project *domain.Project) error {
	err := ps.db.First(project, "path = ?", project.Path).Error
	if err == nil {
		return &ProjectExistsError{Path: project.Path}
	}
	if err == gorm.ErrRecordNotFound {
		return ps.db.Create(project).Error
	}
	return err
}

func (ps *projectStore) LoadProject(path string) (*domain.Project, error) {
	var project domain.Project
	err := ps.db.Preload(clause.Associations).First(&project, "path = ?", path).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &ProjectNotFoundError{Path: path}
	}
	return &project, err
}

func (ps *projectStore) LoadProjectList(paths ...string) ([]*domain.Project, error) {
	numPaths := len(paths)
	if numPaths < 1 {
		return ps.loadAllProjects()
	}
	projects := make([]*domain.Project, 0, numPaths)
	for _, path := range paths {
		project, err := ps.LoadProject(path)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (ps *projectStore) loadAllProjects() ([]*domain.Project, error) {
	var projects []*domain.Project
	err := ps.db.Preload(clause.Associations).Find(&projects).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &NoProjectsFoundError{}
	}
	return projects, err
}

func (ps *projectStore) UpdateProjectLocation(oldPath, newPath string) error {
	err := ps.db.Model(&domain.Project{Path: oldPath}).Update("path", newPath).Error
	if err == gorm.ErrRecordNotFound {
		return &ProjectNotFoundError{Path: oldPath}
	}
	return err
}

func (ps *projectStore) RemoveProject(path string) error {
	err := ps.db.Delete(&domain.Project{}, "path = ? AND deleted_at IS NULL", path).Error
	if err == gorm.ErrRecordNotFound {
		return &ProjectNotFoundError{Path: path}
	}
	return err
}

func (ps *projectStore) PurgeProject(path string) error {
	err := ps.db.Unscoped().Delete(&domain.Project{}, "path = ?", path).Error
	if err == gorm.ErrRecordNotFound {
		return &ProjectNotFoundError{Path: path}
	}
	return err
}
