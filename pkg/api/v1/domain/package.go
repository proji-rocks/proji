package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type (
	// DirEntry represents a directory entry. This is used to represent a directory or file in a DirTree
	DirEntry struct {
		Path     string    `json:"path" toml:"path"`                             // Path is the path of the entry in a project
		IsDir    bool      `json:"is_dir" toml:"is_dir"`                         // IsDir indicates if the entry is a directory
		Template *Template `json:"template,omitempty" toml:"template,omitempty"` // Template is an optional file that will be rendered instead of an empty file
	}

	// DirEntryConfig represents a directory entry configuration. It is used as part of the DirTreeConfig struct.
	DirEntryConfig struct {
		Path     string          `json:"path" toml:"path"`
		IsDir    bool            `json:"is_dir" toml:"is_dir"`
		Template *TemplateConfig `json:"template,omitempty" toml:"template,omitempty"`
	}

	// DirTree represents a directory tree. This is used to represent a directory tree in a package.
	//
	// Note: Using the toml tag 'entry' here instead of 'entries' is intentional. This is to avoid grammatically
	// incorrect sentences when using the toml package. For example, if we used 'entries' here, the toml Package
	// would render a single entry as 'dir_tree.entries = ...' which is grammatically incorrect. Using 'entry'
	// here makes it render as 'dir_tree.entry = ...' which is grammatically correct.
	DirTree struct {
		Entries []*DirEntry `json:"entries" toml:"entry"`
	}

	// DirTreeConfig represents a directory tree configuration. It is used as part of the PackageConfig struct.
	DirTreeConfig struct {
		Entries []*DirEntryConfig `json:"entries" toml:"entry"`
	}

	// Package represents a package. Package is meant to be used for display purposes as it loads all info about a
	// package that might of interest to the user. It is not meant to be used for storage purposes.
	Package struct {
		Label       string           `json:"label" toml:"label"`
		Name        string           `json:"name" toml:"name"`
		UpstreamURL *string          `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty" toml:"sha,omitempty"`
		Description *string          `json:"description,omitempty" toml:"description,omitempty"`
		DirTree     *DirTree         `json:"dir_tree,omitempty" toml:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty" toml:"plugins,omitempty"`
		CreatedAt   time.Time        `json:"created_at" toml:"created_at"`
		UpdatedAt   time.Time        `json:"updated_at" toml:"updated_at"`
	}

	// PackageConfig represents a package configuration. PackageConfig is meant to be used for storage purposes as
	// it only loads the bare minimum of info about a package that is needed to import and export the package. It is not
	// meant to be used for display purposes.
	// Note; Template and Plugin also need to be made suitable for storage. They contain ID, created_at and updated_at
	// fields. These fields are not needed for storage purposes.
	PackageConfig struct {
		Label       string                 `json:"label" toml:"label"`
		Name        string                 `json:"name" toml:"name"`
		UpstreamURL *string                `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		SHA         *string                `json:"sha,omitempty" toml:"sha,omitempty"`
		Description *string                `json:"description,omitempty" toml:"description,omitempty"`
		DirTree     *DirTreeConfig         `json:"dir_tree,omitempty" toml:"dir_tree,omitempty"`
		Plugins     *PluginSchedulerConfig `json:"plugins,omitempty" toml:"plugins,omitempty"`
	}

	// PackageAdd is used to add new packages to the database.
	PackageAdd struct {
		Label       string           `json:"label" toml:"label"`
		Name        string           `json:"name" toml:"name"`
		UpstreamURL *string          `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty" toml:"sha,omitempty"`
		Description *string          `json:"description,omitempty" toml:"description,omitempty"`
		DirTree     *DirTree         `json:"dir_tree,omitempty" toml:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty" toml:"plugins,omitempty"`
	}

	// PackageUpdate is used to update packages in the database.
	PackageUpdate struct {
		Label       string           `json:"label" toml:"label"`
		Name        string           `json:"name,omitempty" toml:"name,omitempty"`
		UpstreamURL *string          `json:"upstream_url,omitempty" toml:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty" toml:"sha,omitempty"`
		Description *string          `json:"description,omitempty" toml:"description,omitempty"`
		DirTree     *DirTree         `json:"dir_tree,omitempty" toml:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty" toml:"plugins,omitempty"`
	}

	// PackageService is used to manage packages, typically by calling a PackageRepo under the hood.
	PackageService interface {
		Fetch(ctx context.Context) ([]Package, error)
		GetByLabel(ctx context.Context, label string) (Package, error)
		Store(ctx context.Context, _package *PackageAdd) error
		Update(ctx context.Context, _package *PackageUpdate) error
		UpdateFromUpstream(ctx context.Context, _package *PackageUpdate) error
		Remove(ctx context.Context, label string) error
	}

	// PackageRepo is used to fetch packages from the database.
	PackageRepo interface {
		Fetch(ctx context.Context) ([]Package, error)
		GetByLabel(ctx context.Context, label string) (Package, error)
		Store(ctx context.Context, _package *PackageAdd) error
		Update(ctx context.Context, _package *PackageUpdate) error
		Remove(ctx context.Context, label string) error
	}
)

const bucketPackages = "packages"

// Bucket returns the bucket name for the package.
func (*Package) Bucket() string {
	return bucketPackages
}

// MarshalJSON marshals the package into JSON. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (p *PackageAdd) MarshalJSON() ([]byte, error) {
	type Alias PackageAdd
	now := time.Now().UTC

	return json.Marshal(&struct {
		*Alias
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: now(),
		UpdatedAt: now(),
	})
}

// MarshalTOML marshals the package into TOML. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (p *PackageAdd) MarshalTOML() ([]byte, error) {
	type Alias PackageAdd
	now := time.Now().UTC

	return toml.Marshal(&struct {
		*Alias
		CreatedAt time.Time `toml:"created_at"`
		UpdatedAt time.Time `toml:"updated_at"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: now(),
		UpdatedAt: now(),
	})
}

// AsUpdatable converts a package to an updatable package. Typically, this is used to edit or update an already
// loaded package.
func (p *Package) AsUpdatable() *PackageUpdate {
	return &PackageUpdate{
		Label:       p.Label,
		Name:        p.Name,
		UpstreamURL: p.UpstreamURL,
		SHA:         p.SHA,
		Description: p.Description,
		DirTree:     p.DirTree,
		Plugins:     p.Plugins,
	}
}

// NewPackage creates a new package with the given name and label.
func NewPackage(name, label string) *PackageAdd {
	if name == "" {
		name = "Unknown"
	}
	if label == "" {
		label = "xxx"
	}

	return &PackageAdd{Name: name, Label: label}
}

// NewPackageWithAutoLabel creates a new package with the given name. The label gets auto-generated based on the name.
func NewPackageWithAutoLabel(name string) *PackageAdd {
	return NewPackage(name, generateLabelFromName(name))
}

func (e *DirEntry) toConfig() *DirEntryConfig {
	if e == nil {
		return nil
	}

	return &DirEntryConfig{
		Path:     e.Path,
		IsDir:    e.IsDir,
		Template: e.Template.toConfig(),
	}
}

func (d *DirTree) toConfig() *DirTreeConfig {
	if d == nil {
		return nil
	}

	conf := &DirTreeConfig{
		Entries: make([]*DirEntryConfig, len(d.Entries)),
	}

	for i, entry := range d.Entries {
		conf.Entries[i] = entry.toConfig()
	}

	return conf
}

func (p *Package) ToConfig() *PackageConfig {
	return &PackageConfig{
		Label:       p.Label,
		Name:        p.Name,
		UpstreamURL: p.UpstreamURL,
		SHA:         p.SHA,
		Description: p.Description,
		DirTree:     p.DirTree.toConfig(),
		Plugins:     p.Plugins.toConfig(),
	}
}
