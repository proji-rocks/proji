package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"
)

type (
	// Project represents a package. Project is meant to be used for display purposes as it loads all info about a
	// package that might of interest to the user. It is not meant to be used for storage purposes.
	Project struct {
		ID          string    `json:"id"`
		Path        string    `json:"path"`
		Name        string    `json:"name"`
		Package     string    `json:"package"`
		Description *string   `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// ProjectAdd is used to add new packages to the database.
	ProjectAdd struct {
		Path        string  `json:"path"`
		Name        string  `json:"name"`
		Package     string  `json:"package"`
		Description *string `json:"description,omitempty"`
	}

	// ProjectUpdate is used to update packages in the database.
	ProjectUpdate struct {
		ID          string  `json:"id"`
		Path        string  `json:"path"`
		Name        string  `json:"name,omitempty"`
		Package     string  `json:"package,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// ProjectService is used to manage packages, typically by calling a ProjectRepo under the hood.
	ProjectService interface {
		Fetch(ctx context.Context) ([]Project, error)
		GetByID(ctx context.Context, id string) (Project, error)
		Store(ctx context.Context, project *ProjectAdd) error
		Update(ctx context.Context, project *ProjectUpdate) error
		Remove(ctx context.Context, label string) error
	}

	// ProjectRepo is used to fetch packages from the database.
	ProjectRepo interface {
		ProjectService
	}
)

const bucketProjects = "projects"

// Bucket returns the bucket name for the package.
func (Project) Bucket() string {
	return bucketProjects
}

// MarshalJSON marshals the project into JSON. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (p *ProjectAdd) MarshalJSON() ([]byte, error) {
	type Alias ProjectAdd

	// ID should be applied only if it is not set. We overwrite the ID field with a new ID if it is not set so that
	// the marshaled JSON contains the new ID as well as the receiving project instance.
	/*	if p.ID == "" {
			p.ID = xid.New().String()
		}
	*/
	return json.Marshal(&struct {
		*Alias
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Alias:     (*Alias)(p),
		ID:        xid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// NewProject creates a new package with the given name and label.
func NewProject(packageLabel, path, name string) *ProjectAdd {
	return &ProjectAdd{
		Path:    path,
		Name:    name,
		Package: packageLabel,
	}
}
