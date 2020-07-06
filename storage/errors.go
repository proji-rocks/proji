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

// ClassNotFoundError represents an error for the case that a query for a class returns a
// gorm.ErrRecordNotFound error.
type ClassNotFoundError struct {
	Label string
}

//
func (e *ClassNotFoundError) Error() string {
	return fmt.Sprintf("class with label '%s' not found", e.Label)
}

// NewClassNotFoundError returns a pointer to an initialized ClassNotFoundError object.
func NewClassNotFoundError(label string) *ClassNotFoundError {
	return &ClassNotFoundError{Label: label}
}

// NoClassesFoundError represents an error for the case that no classes were found by a query.
type NoClassesFoundError struct{}

func (e *NoClassesFoundError) Error() string {
	return "no classes were found"
}

// NewNoClassesFoundError returns a pointer to an initialized NoClassesFoundError object.
func NewNoClassesFoundError() *NoClassesFoundError {
	return &NoClassesFoundError{}
}

// ClassExistsError represents an error for the case that a query for a class returns no result.
type ClassExistsError struct {
	Label string
}

func (e *ClassExistsError) Error() string {
	return fmt.Sprintf("class with label '%s' already exists", e.Label)
}

// NewClassExistsError returns a pointer to an initialized ClassExistsError object.
func NewClassExistsError(label string) *ClassExistsError {
	return &ClassExistsError{Label: label}
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
