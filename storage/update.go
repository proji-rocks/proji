package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
)

type UpdateService interface {
	UpdateProjectLocation(oldPath, newPath string) error // UpdateProjectLocation updates the path of a project in storage.
}

// UpdateProjectLocation updates the location of a project in storage.
func (db *Database) UpdateProjectLocation(oldPath, newPath string) error {
	err := db.Connection.Model(&models.Project{Path: oldPath}).Update("path", newPath).Error
	if err == gorm.ErrRecordNotFound {
		return &ProjectNotFoundError{Path: oldPath}
	}
	return err
}
