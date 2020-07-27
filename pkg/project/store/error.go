package projectstore

import (
	"errors"
)

// ErrProjectNotFound represents an error for the case that a query for a project returns a
// gorm.ErrRecordNotFound error.
var ErrProjectNotFound = errors.New("project not found")

// ErrNoProjectsFound represents an error for the case that no projects were found by a query.
var ErrNoProjectsFound = errors.New("no projects found")

// ErrProjectExists represents an error for the case that a query for a project returns no result.
var ErrProjectExists = errors.New("a project is already assigned to the path")
