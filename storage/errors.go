package storage

import "fmt"

// UnsupportedDatabaseDialectError represents an error for the case that the user passed a db driver for a
// unsupported database dialect. Proji uses Gorm under the hood so take a look at its docs for a list of
// documented dialects.
// https://gorm.io/docs/connecting_to_the_database.html#Supported-Databases
type UnsupportedDatabaseDialectError struct {
	Dialect string
}

func (e *UnsupportedDatabaseDialectError) Error() string {
	return fmt.Sprintf("%s is not in the list of supported database dialects", e.Dialect)
}

// PackageNotFoundError represents an error for the case that a query for a package returns a
// gorm.ErrRecordNotFound error.
type PackageNotFoundError struct {
	Label string
}

func (e *PackageNotFoundError) Error() string {
	return fmt.Sprintf("package with label '%s' not found", e.Label)
}

// NoPackagesFoundError represents an error for the case that no packages were found by a query.
type NoPackagesFoundError struct{}

func (e *NoPackagesFoundError) Error() string {
	return "no packages were found"
}

// PackageExistsError represents an error for the case that a query for a package returns no result.
type PackageExistsError struct {
	Label string
}

func (e *PackageExistsError) Error() string {
	return fmt.Sprintf("package with label '%s' already exists", e.Label)
}

// ProjectNotFoundError represents an error for the case that a query for a project returns a
// gorm.ErrRecordNotFound error.
type ProjectNotFoundError struct {
	Path string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project at path '%s' not found", e.Path)
}

// NoProjectsFoundError represents an error for the case that no projects were found by a query.
type NoProjectsFoundError struct{}

func (e *NoProjectsFoundError) Error() string {
	return "no projects were found"
}

// ProjectExistsError represents an error for the case that a query for a project returns no result.
type ProjectExistsError struct {
	Path string
}

func (e *ProjectExistsError) Error() string {
	return fmt.Sprintf("a project is already assigned to the path '%s'", e.Path)
}
