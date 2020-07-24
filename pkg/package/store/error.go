package packagestore

import (
	"errors"
)

// ErrPackageNotFound represents an error for the case that a query for a package returns a
// gorm.ErrRecordNotFound error.
var ErrPackageNotFound = errors.New("package not found")

// ErrPackageExists represents an error for the case that a query for a package returns no result.
var ErrPackageExists = errors.New("package already exists")
