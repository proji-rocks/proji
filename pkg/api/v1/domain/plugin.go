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
		ID          string    `json:"id" toml:"id"`
		Path        string    `json:"path" toml:"path"`
		UpstreamURL *string   `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		Description *string   `json:"description,omitempty" toml:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at" toml:"created_at"`
		UpdatedAt   time.Time `json:"updated_at" toml:"updated_at"`
	}

	// PluginConfig represents a plugin configuration. It is used as part of the PackageConfig.
	PluginConfig struct {
		Path        string  `json:"path" toml:"path"`
		UpstreamURL *string `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty" toml:"description,omitempty"`
	}

	// PluginScheduler is used to schedule plugins. It has two lists of plugins: one for the pre-creation and one for
	// the post-creation of a new project. The pre-creation list is executed before the project is created and the
	// post-creation list is executed after the project is created. The scheduler follows the order of the lists.
	PluginScheduler struct {
		Pre  []*Plugin `json:"pre,omitempty" toml:"pre,omitempty"`   // Pre-creation plugins.
		Post []*Plugin `json:"post,omitempty" toml:"post,omitempty"` // Post-creation plugins.
	}

	// PluginSchedulerConfig represents a plugin scheduler configuration. It is used as part of the PackageConfig.
	PluginSchedulerConfig struct {
		Pre  []*PluginConfig `json:"pre,omitempty" toml:"pre,omitempty"`   // Pre-creation plugins.
		Post []*PluginConfig `json:"post,omitempty" toml:"post,omitempty"` // Post-creation plugins.
	}

	// PluginAdd represents a project to be added.
	PluginAdd struct {
		Path        string  `json:"path" toml:"path"`
		UpstreamURL *string `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty" toml:"description,omitempty"`
	}

	// PluginUpdate represents a project to be updated.
	PluginUpdate struct {
		ID          string  `json:"id" toml:"id"`
		Path        string  `json:"path" toml:"path"`
		UpstreamURL *string `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		Description *string `json:"description,omitempty" toml:"description,omitempty"`
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
func (*Plugin) Bucket() string {
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

func (p *Plugin) toConfig() *PluginConfig {
	if p == nil {
		return nil
	}

	return &PluginConfig{
		Path:        p.Path,
		UpstreamURL: p.UpstreamURL,
		Description: p.Description,
	}
}

func (p *PluginScheduler) toConfig() *PluginSchedulerConfig {
	if p == nil {
		return nil
	}

	conf := &PluginSchedulerConfig{
		Pre:  make([]*PluginConfig, len(p.Pre)),
		Post: make([]*PluginConfig, len(p.Post)),
	}

	for _, plg := range p.Pre {
		conf.Pre = append(conf.Pre, plg.toConfig())
	}

	for _, plg := range p.Post {
		conf.Post = append(conf.Post, plg.toConfig())
	}

	return conf
}
