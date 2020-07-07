package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoadService interface {
	LoadClass(label string) (*models.Class, error)           // LoadClass loads a class from storage by its label.
	LoadClasses(labels ...string) ([]*models.Class, error)   // LoadClasses returns classes by the given labels. If no labels are given, all classes are loaded.
	LoadProject(path string) (*models.Project, error)        // LoadProject loads a project from storage by its path.
	LoadProjects(paths ...string) ([]*models.Project, error) // LoadProjects returns projects by the given paths. If no paths are given, all projects are loaded.
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

// LoadClasses loads classes by the given labels. If not labels are given, all classes are loaded.
func (db *Database) LoadClasses(labels ...string) ([]*models.Class, error) {
	lenLabels := len(labels)
	if lenLabels < 1 {
		return db.loadAllClasses()
	}
	classes := make([]*models.Class, 0, lenLabels)
	for _, label := range labels {
		class, err := db.LoadClass(label)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, nil
}

// loadAllClasses loads and returns all classes found in the database.
func (db *Database) loadAllClasses() ([]*models.Class, error) {
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

// LoadProjects returns projects by the given paths. If no paths are given, all projects are loaded.
func (db *Database) LoadProjects(paths ...string) ([]*models.Project, error) {
	numPaths := len(paths)
	if numPaths < 1 {
		return db.loadAllProjects()
	}
	projects := make([]*models.Project, 0, numPaths)
	for _, path := range paths {
		project, err := db.LoadProject(path)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// loadAllProjects loads and returns all projects found in the database.
func (db *Database) loadAllProjects() ([]*models.Project, error) {
	var projects []*models.Project
	err := db.Connection.Preload(clause.Associations).Find(&projects).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewNoProjectsFoundError()
	}
	return projects, err
}
