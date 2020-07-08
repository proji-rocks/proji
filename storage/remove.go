package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
)

type RemoveService interface {
	RemovePackage(label string) error // RemovePackage removes a package from storage.
	PurgePackage(label string) error  // PurgePackage removes a soft-deleted package finally from storage.
	RemoveProject(path string) error  // RemoveProject removes a project from storage.
	PurgeProject(path string) error   // PurgeProject removes a soft-deleted project finally from storage.
}

// RemovePackage performs a soft-delete of a given package from storage.
func (db *Database) RemovePackage(label string) error {
	err := db.Connection.Delete(&models.Package{}, "label = ? AND deleted_at IS NULL", label).Error
	if err == gorm.ErrRecordNotFound {
		return NewPackageNotFoundError(label)
	}
	return err
}

// PurgePackage removes a soft-deleted package finally from storage.
func (db *Database) PurgePackage(label string) error {
	err := db.Connection.Unscoped().Delete(&models.Package{}, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return NewPackageNotFoundError(label)
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
