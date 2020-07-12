package projectstore

import "fmt"

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
