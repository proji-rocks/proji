package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoadService interface {
	LoadClass(label string) (*models.Class, error)    // LoadClass loads a class from storage by its label.
	LoadAllClasses() ([]*models.Class, error)         // LoadAllClasses loads all available classes from storage.
	LoadProject(path string) (*models.Project, error) // LoadProject loads a project from storage by its path.
	LoadAllProjects() ([]*models.Project, error)      // LoadAllProjects returns a list of all projects in storage.
}

// LoadClass loads a class from storage by its label.
func (db *Database) LoadClass(label string) (*models.Class, error) {
	var class models.Class
	err := db.Connection.Preload(clause.Associations).First(&class, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewClassNotFoundError(label)
	}
	return &class, err
}

// LoadAllClasses loads all available classes from storage.
func (db *Database) LoadAllClasses() ([]*models.Class, error) {
	var classes []*models.Class
	err := db.Connection.Preload(clause.Associations).Find(&classes).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewNoClassesFoundError()
	}
	return classes, err
}

// LoadProject loads a project from storage by its path.
func (db *Database) LoadProject(path string) (*models.Project, error) {
	var project models.Project
	err := db.Connection.Preload(clause.Associations).First(&project, "path = ?", path).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewProjectNotFoundError(path)
	}
	return &project, err
}

// LoadAllProjects returns a list of all projects in storage.
func (db *Database) LoadAllProjects() ([]*models.Project, error) {
	var projects []*models.Project
	err := db.Connection.Preload(clause.Associations).Find(&projects).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewNoProjectsFoundError()
	}
	return projects, err
}
