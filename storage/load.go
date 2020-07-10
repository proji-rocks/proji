package storage

import (
	"github.com/nikoksr/proji/storage/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoadService interface {
	LoadPackage(label string) (*models.Package, error)        // LoadPackage loads a package from storage by its label.
	LoadPackages(labels ...string) ([]*models.Package, error) // LoadPackages returns packages by the given labels. If no labels are given, all packages are loaded.
	LoadProject(path string) (*models.Project, error)         // LoadProject loads a project from storage by its path.
	LoadProjects(paths ...string) ([]*models.Project, error)  // LoadProjects returns projects by the given paths. If no paths are given, all projects are loaded.
}

// LoadPackage loads a package from storage by its label.
func (db *Database) LoadPackage(label string) (*models.Package, error) {
	var pkg models.Package
	err := db.Connection.Preload(clause.Associations).First(&pkg, "label = ?", label).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &PackageNotFoundError{Label: label}
	}
	return &pkg, err
}

// LoadPackages loads packages by the given labels. If not labels are given, all packages are loaded.
func (db *Database) LoadPackages(labels ...string) ([]*models.Package, error) {
	lenLabels := len(labels)
	if lenLabels < 1 {
		return db.loadAllPackages()
	}
	packages := make([]*models.Package, 0, lenLabels)
	for _, label := range labels {
		pkg, err := db.LoadPackage(label)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// loadAllPackages loads and returns all packages found in the database.
func (db *Database) loadAllPackages() ([]*models.Package, error) {
	var packages []*models.Package
	err := db.Connection.Preload(clause.Associations).Find(&packages).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &NoPackagesFoundError{}
	}
	return packages, err
}

// LoadProject loads a project from storage by its path.
func (db *Database) LoadProject(path string) (*models.Project, error) {
	var project models.Project
	err := db.Connection.Preload(clause.Associations).First(&project, "path = ?", path).Error
	if err == gorm.ErrRecordNotFound {
		return nil, &ProjectNotFoundError{Path: path}
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
		return nil, &NoProjectsFoundError{}
	}
	return projects, err
}
