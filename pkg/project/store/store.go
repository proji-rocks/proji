package projectstore

import (
	"github.com/nikoksr/proji/pkg/domain"
	"github.com/pkg/errors"
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
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	err = tx.Omit(clause.Associations).Create(project).Error
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert project")
	}

	return tx.Commit().Error
}

func (ps *projectStore) LoadProject(path string) (*domain.Project, error) {
	var project domain.Project
	tx := ps.db.Where("path = ?", path).First(&project)
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
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(&domain.Project{Path: oldPath}).Update("path", newPath).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || tx.RowsAffected < 1 {
		tx.Rollback()
		return ErrProjectNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ps *projectStore) RemoveProject(path string) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	err := tx.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Where("path = ?", path).Delete(&domain.Project{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return ErrProjectNotFound
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
