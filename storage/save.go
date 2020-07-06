package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
)

type SaveService interface {
	SaveClass(class *models.Class) error       // SaveClass saves a class to storage.
	SaveProject(project *models.Project) error // SaveProject saves a project to storage.
}

// SaveClass saves a class to storage.
func (db *Database) SaveClass(class *models.Class) error {
	err := db.Connection.First(class).Error
	if err == nil {
		return NewClassExistsError(class.Label, class.Name)
	}
	if err == gorm.ErrRecordNotFound {
		return db.Connection.Create(class).Error
	}
	return err
}

// SaveProject saves a project to storage.
func (db *Database) SaveProject(project *models.Project) error {
	err := db.Connection.First(project).Error
	if err == nil {
		return NewProjectExistsError(project.Path)
	}
	if err == gorm.ErrRecordNotFound {
		return db.Connection.Create(project).Error
	}
	return err
}
