package storage

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/storage/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *Database) Migrate() error {
	modelList := []interface{}{
		&models.Class{},
		&models.Plugin{},
		&models.Project{},
		&models.Template{},
	}
	for _, model := range modelList {
		err := db.Connection.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to auto-migrate model, %s", err.Error())
		}
	}
	return nil
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

// UpdateProjectLocation updates the location of a project in storage.
func (db *Database) UpdateProjectLocation(oldPath, newPath string) error {
	err := db.Connection.Model(&models.Project{Path: oldPath}).Update("path", newPath).Error
	if err == gorm.ErrRecordNotFound {
		return NewProjectNotFoundError(oldPath)
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
