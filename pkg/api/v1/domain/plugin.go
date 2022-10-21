package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"
)

type (
	// Plugin represents a package project. Plugins are usually some kind of scripts (currently only lua) that
	// are executed by the package manager.
	Plugin struct {
		ID          string    `json:"id"`                     // ID is the unique identifier of the project
		Path        string    `json:"path"`                   // Path is the path to the project
		UpstreamURL *string   `json:"upstream_url,omitempty"` // UpstreamURL is the URL of the upstream project
		Description *string   `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// PluginScheduler is used to schedule plugins. It has two lists of plugins: one for the pre-creation and one for
	// the post-creation of a new project. The pre-creation list is executed before the project is created and the
	// post-creation list is executed after the project is created. The scheduler follows the order of the lists.
	PluginScheduler struct {
		Pre  []*Plugin `json:"pre,omitempty"`  // List of plugins to run before the build
		Post []*Plugin `json:"post,omitempty"` // List of plugins to run after the build
	}

	// PluginAdd represents a project to be added.
	PluginAdd struct {
		Name        string  `json:"name"`
		Path        string  `json:"path"`
		UpstreamURL *string `json:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// PluginUpdate represents a project to be updated.
	PluginUpdate struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Path        string  `json:"path"`
		UpstreamURL *string `json:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	// PluginService is used to manage plugins, typically by calling a PluginRepo under the hood.
	PluginService interface {
		Fetch(ctx context.Context) ([]Plugin, error)
		GetByID(ctx context.Context, id string) (Plugin, error)
		Store(ctx context.Context, plg *PluginAdd) error
		Update(ctx context.Context, plg *PluginUpdate) error
		Remove(ctx context.Context, id string) error
	}

	// PluginRepo is used to fetch plugins from the database.
	PluginRepo interface {
		PluginService
	}
)

const bucketPlugins = "plugins"

// Bucket returns the bucket name for the project.
func (Plugin) Bucket() string {
	return bucketPlugins
}

// MarshalJSON marshals the project into JSON. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (p *PluginAdd) MarshalJSON() ([]byte, error) {
	type Alias PluginAdd

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
