package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"
)

type (
	// Template represents a package template. Templates are used to create bootstrapped files in projects.
	Template struct {
		ID          string    `json:"id"`                     // ID is the unique identifier of the template
		Path        string    `json:"path"`                   // Path is the path to the template
		UpstreamURL *string   `json:"upstream_url,omitempty"` // UpstreamURL is the URL of the upstream template
		Description *string   `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// TemplateAdd is used to add a new template.
	TemplateAdd struct {
		Name        string  `json:"name"`
		Path        string  `json:"path"`
		Destination string  `json:"destination"`
		IsFile      bool    `json:"is_file"`
		UpstreamURL *string `json:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// TemplateUpdate is used to update an existing template.
	TemplateUpdate struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Path        string  `json:"path"`
		Destination string  `json:"destination"`
		IsFile      bool    `json:"is_file"`
		UpstreamURL *string `json:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// TemplateService is used to manage templates, typically by calling a TemplateRepo under the hood.
	TemplateService interface {
		Store(ctx context.Context, tmpl *TemplateAdd) error
		GetByID(ctx context.Context, id string) (Template, error)
		Fetch(ctx context.Context) ([]Template, error)
		Update(ctx context.Context, tmpl *TemplateUpdate) error
		Remove(ctx context.Context, id string) error
	}

	// TemplateRepo is used to fetch templates from the database.
	TemplateRepo interface {
		TemplateService
	}
)

const bucketTemplates = "templates"

// Bucket returns the bucket name for the template.
func (Template) Bucket() string {
	return bucketTemplates
}

// MarshalJSON marshals the template into JSON. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (t *TemplateAdd) MarshalJSON() ([]byte, error) {
	type Alias TemplateAdd

	return json.Marshal(&struct {
		*Alias
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Alias:     (*Alias)(t),
		ID:        xid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}
