package storage

import (
	"fmt"
	"strings"

	"github.com/nikoksr/proji/pkg/storage/models"
)

// Service interface describes the behaviour of a storage service.
type Service interface {
	Migrate() error                                      // Migrate models to storage.
	SaveClass(class *models.Class) error                 // SaveClass saves a class to storage.
	LoadClass(label string) (*models.Class, error)       // LoadClass loads a class from storage by its label.
	LoadAllClasses() ([]*models.Class, error)            // LoadAllClasses loads all available classes from storage.
	RemoveClass(label string) error                      // RemoveClass removes a class from storage.
	PurgeClass(label string) error                       // PurgeClass removes a soft-deleted class finally from storage.
	SaveProject(project *models.Project) error           // SaveProject saves a project to storage.
	LoadProject(path string) (*models.Project, error)    // LoadProject loads a project from storage by its path.
	LoadAllProjects() ([]*models.Project, error)         // LoadAllProjects returns a list of all projects in storage.
	UpdateProjectLocation(oldPath, newPath string) error // UpdateProjectLocation updates the path of a project in storage.
	RemoveProject(path string) error                     // RemoveProject removes a project from storage.
	PurgeProject(path string) error                      // PurgeProject removes a soft-deleted project finally from storage.
}

// NewService returns a new storage service interface initialized with a given storage driver and connection string.
func NewService(driver, connectionString string) (Service, error) {
	var err error
	driver, connectionString, err = validateParameters(driver, connectionString)
	if err != nil {
		return nil, err
	}

	var svc Service
	if isDatabaseDriver(driver) {
		svc, err = newDatabaseService(driver, connectionString)
	} else {
		return nil, fmt.Errorf("storage service driver %s is not supported", driver)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create storage service, %s", err.Error())
	}

	err = svc.Migrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate models into storage service, %s", err.Error())
	}

	return svc, nil
}

// validateParameters validates that the parameters for driver and connectionString are valid. Driver gets replaced
// by the default driver if it was not given. ConnectionString may not be empty.
func validateParameters(driver, connectionString string) (string, string, error) {
	// Normalize parameters
	driver = strings.TrimSpace(driver)
	connectionString = strings.TrimSpace(connectionString)

	// Check if anything was given. In case of that no driver was given use default driver as fallback value.
	// Error when no connection string was given.
	if len(driver) < 1 {
		driver = defaultDriver
	}
	if len(connectionString) < 0 {
		return "", "", fmt.Errorf("storage service connection string may not be empty")
	}
	return driver, connectionString, nil
}
