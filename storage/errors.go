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

// NewUnsupportedDatabaseDialectError returns a pointer to an initialized UnsupportedDatabaseDialectError object.
func NewUnsupportedDatabaseDialectError(dialect string) *UnsupportedDatabaseDialectError {
	return &UnsupportedDatabaseDialectError{Dialect: dialect}
}

// PackageNotFoundError represents an error for the case that a query for a package returns a
// gorm.ErrRecordNotFound error.
type PackageNotFoundError struct {
	Label string
}

//
func (e *PackageNotFoundError) Error() string {
	return fmt.Sprintf("package with label '%s' not found", e.Label)
}

// NewPackageNotFoundError returns a pointer to an initialized PackageNotFoundError object.
func NewPackageNotFoundError(label string) *PackageNotFoundError {
	return &PackageNotFoundError{Label: label}
}

// NoPackagesFoundError represents an error for the case that no packages were found by a query.
type NoPackagesFoundError struct{}

func (e *NoPackagesFoundError) Error() string {
	return "no packages were found"
}

// NewNoPackagesFoundError returns a pointer to an initialized NoPackagesFoundError object.
func NewNoPackagesFoundError() *NoPackagesFoundError {
	return &NoPackagesFoundError{}
}

// PackageExistsError represents an error for the case that a query for a package returns no result.
type PackageExistsError struct {
	Label string
}

func (e *PackageExistsError) Error() string {
	return fmt.Sprintf("package with label '%s' already exists", e.Label)
}

// NewPackageExistsError returns a pointer to an initialized PackageExistsError object.
func NewPackageExistsError(label string) *PackageExistsError {
	return &PackageExistsError{Label: label}
}

// ProjectNotFoundError represents an error for the case that a query for a project returns a
// gorm.ErrRecordNotFound error.
type ProjectNotFoundError struct {
	Path string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project at path '%s' not found", e.Path)
}

// NewProjectNotFoundError returns a pointer to an initialized ProjectNotFoundError object.
func NewProjectNotFoundError(path string) *ProjectNotFoundError {
	return &ProjectNotFoundError{Path: path}
}

// NoProjectsFoundError represents an error for the case that no projects were found by a query.
type NoProjectsFoundError struct{}

func (e *NoProjectsFoundError) Error() string {
	return "no projects were found"
}

// NewNoProjectsFoundError returns a pointer to an initialized NoProjectsFoundError object.
func NewNoProjectsFoundError() *NoProjectsFoundError {
	return &NoProjectsFoundError{}
}

// ProjectExistsError represents an error for the case that a query for a project returns no result.
type ProjectExistsError struct {
	Path string
}

func (e *ProjectExistsError) Error() string {
	return fmt.Sprintf("a project is already assigned to the path '%s'", e.Path)
}

// NewProjectExistsError returns a pointer to an initialized ProjectExistsError object.
func NewProjectExistsError(path string) *ProjectExistsError {
	return &ProjectExistsError{Path: path}
}
