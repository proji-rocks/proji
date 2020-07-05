package storage

import "fmt"

// ClassNotFoundError represents an error for the case that a query for a class returns a
// gorm.ErrRecordNotFound error.
type ClassNotFoundError struct {
	Label string
}

//
func (e *ClassNotFoundError) Error() string {
	return fmt.Sprintf("class with label %s not found\n", e.Label)
}

// NewClassNotFoundError returns a pointer to an initialized ClassNotFoundError object.
func NewClassNotFoundError(label string) *ClassNotFoundError {
	return &ClassNotFoundError{Label: label}
}

// NoClassesFoundError represents an error for the case that no classes were found by a query.
type NoClassesFoundError struct{}

func (e *NoClassesFoundError) Error() string {
	return fmt.Sprintf("no classes were found\n")
}

// NewNoClassesFoundError returns a pointer to an initialized NoClassesFoundError object.
func NewNoClassesFoundError() *NoClassesFoundError {
	return &NoClassesFoundError{}
}

// ClassExistsError represents an error for the case that a query for a class returns no result.
type ClassExistsError struct {
	Label string
	Name  string
}

func (e *ClassExistsError) Error() string {
	return fmt.Sprintf("class %s(%s) already exists\n", e.Name, e.Label)
}

// NewClassExistsError returns a pointer to an initialized ClassExistsError object.
func NewClassExistsError(label, name string) *ClassExistsError {
	return &ClassExistsError{Label: label, Name: name}
}

// ProjectNotFoundError represents an error for the case that a query for a project returns a
// gorm.ErrRecordNotFound error.
type ProjectNotFoundError struct {
	Path string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project at path %s not found\n", e.Path)
}

// NewProjectNotFoundError returns a pointer to an initialized ProjectNotFoundError object.
func NewProjectNotFoundError(path string) *ProjectNotFoundError {
	return &ProjectNotFoundError{Path: path}
}

// NoProjectsFoundError represents an error for the case that no projects were found by a query.
type NoProjectsFoundError struct{}

func (e *NoProjectsFoundError) Error() string {
	return fmt.Sprintf("no projects were found\n")
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
	return fmt.Sprintf("a project is already assigned to the path %s\n", e.Path)
}

// NewProjectExistsError returns a pointer to an initialized ProjectExistsError object.
func NewProjectExistsError(path string) *ProjectExistsError {
	return &ProjectExistsError{Path: path}
}
