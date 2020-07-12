package packagestore

import "fmt"

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
