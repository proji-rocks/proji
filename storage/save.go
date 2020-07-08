package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
)

type SaveService interface {
	SavePackage(pkg *models.Package) error     // SavePackage saves a package to storage.
	SaveProject(project *models.Project) error // SaveProject saves a project to storage.
}

// SavePackage saves a package to storage.
func (db *Database) SavePackage(pkg *models.Package) error {
	err := db.Connection.First(pkg, "label = ?", pkg.Label).Error
	if err == nil {
		return NewPackageExistsError(pkg.Label)
	}
	if err == gorm.ErrRecordNotFound {
		return db.Connection.Create(pkg).Error
	}
	return err
}

// SaveProject saves a project to storage.
func (db *Database) SaveProject(project *models.Project) error {
	err := db.Connection.First(project, "path = ?", project.Path).Error
	if err == nil {
		return NewProjectExistsError(project.Path)
	}
	if err == gorm.ErrRecordNotFound {
		return db.Connection.Create(project).Error
	}
	return err
}
