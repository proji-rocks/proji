package projectstore

import (
	"errors"

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
	err := ps.db.Where("path = ?", project.Path).First(project).Error
	if err == nil {
		return ErrProjectExists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ps.db.Create(project).Error
	}
	return err
}

func (ps *projectStore) LoadProject(path string) (*domain.Project, error) {
	var project domain.Project
	tx := ps.db.Preload(clause.Associations).Where("path = ?", path).First(&project)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		return nil, ErrProjectNotFound
	}
	return &project, tx.Error
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoProjectsFound
	}
	return projects, err
}

func (ps *projectStore) UpdateProjectLocation(oldPath, newPath string) error {
	tx := ps.db.Model(&domain.Project{Path: oldPath}).Update("path", newPath)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		return ErrProjectNotFound
	}
	return tx.Error
}

func (ps *projectStore) RemoveProject(path string) error {
	tx := ps.db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Where("path = ?", path).Delete(&domain.Project{})
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		return ErrProjectNotFound
	}
	return tx.Error
}
