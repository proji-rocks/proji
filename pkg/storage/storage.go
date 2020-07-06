package storage

import (
	"fmt"
	"strings"
)

// Service interface describes the behaviour of a storage service.
type Service interface {
	Migrate() error // Migrate models to storage.
	SaveService     // Service to handle save actions for the storage.
	LoadService     // Service to handle load actions for the storage.
	UpdateService   // Service to handle update actions for the storage.
	RemoveService   // Service to handle remove actions for the storage.
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
