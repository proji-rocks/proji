package storage

import (
	"github.com/nikoksr/proji/pkg/storage/models"
	"gorm.io/gorm"
)

type RemoveService interface {
	RemoveClass(label string) error  // RemoveClass removes a class from storage.
	PurgeClass(label string) error   // PurgeClass removes a soft-deleted class finally from storage.
	RemoveProject(path string) error // RemoveProject removes a project from storage.
	PurgeProject(path string) error  // PurgeProject removes a soft-deleted project finally from storage.
}

// RemoveClass performs a soft-delete of a given class from storage.
func (db *Database) RemoveClass(label string) error {
	err := db.Connection.Delete(&models.Class{}, "label = ? AND deleted_at IS NULL", label).Error
	if err == gorm.ErrRecordNotFound {
		return NewClassNotFoundError(label)
	}
	return err
}

// PurgeClass removes a soft-deleted class finally from storage.
func (db *Database) PurgeClass(label string) error {
	err := db.Connection.Unscoped().Delete(&models.Class{}, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return NewClassNotFoundError(label)
	}
	return err
}

// RemoveProject removes a project from storage.
func (db *Database) RemoveProject(path string) error {
	err := db.Connection.Delete(&models.Project{}, "path = ? AND deleted_at IS NULL", path).Error
	if err == gorm.ErrRecordNotFound {
		return NewProjectNotFoundError(path)
	}
	return err
}

// PurgeProject removes a soft-deleted project finally from storage.
func (db *Database) PurgeProject(path string) error {
	err := db.Connection.Unscoped().Delete(&models.Project{}, "path = ?", path).Error
	if err == gorm.ErrRecordNotFound {
		return NewProjectNotFoundError(path)
	}
	return err
}
