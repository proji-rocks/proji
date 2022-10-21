package domain

import (
	"context"
	"encoding/json"
	"time"
)

type (
	// DirEntry represents a directory entry. This is used to represent a directory or file in a DirTree
	DirEntry struct {
		IsDir    bool      `json:"is_dir"`             // IsDir indicates if the entry is a directory
		Path     string    `json:"path"`               // Path is the path of the entry in a project
		Template *Template `json:"template,omitempty"` // Template is an optional file that will be rendered instead of an empty file
	}

	// DirTree represents a directory tree. This is used to represent a directory tree in a package.
	DirTree = []*DirEntry

	// Package represents a package. Package is meant to be used for display purposes as it loads all info about a
	// package that might of interest to the user. It is not meant to be used for storage purposes.
	Package struct {
		Label       string           `json:"label"`
		Name        string           `json:"name"`
		UpstreamURL *string          `json:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty"`
		Description *string          `json:"description,omitempty"`
		DirTree     DirTree          `json:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty"`
		CreatedAt   time.Time        `json:"created_at"`
		UpdatedAt   time.Time        `json:"updated_at"`
	}

	// PackageAdd is used to add new packages to the database.
	PackageAdd struct {
		Label       string           `json:"label"`
		Name        string           `json:"name"`
		UpstreamURL *string          `json:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty"`
		Description *string          `json:"description,omitempty"`
		DirTree     DirTree          `json:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty"`
	}

	// PackageUpdate is used to update packages in the database.
	PackageUpdate struct {
		Label       string           `json:"label"`
		Name        string           `json:"name,omitempty"`
		UpstreamURL *string          `json:"upstream_url,omitempty"`
		SHA         *string          `json:"sha,omitempty"`
		Description *string          `json:"description,omitempty"`
		DirTree     DirTree          `json:"dir_tree,omitempty"`
		Plugins     *PluginScheduler `json:"plugins,omitempty"`
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
func (Package) Bucket() string {
	return bucketPackages
}

// MarshalJSON marshals the package into JSON. It is used to dynamically add timestamps for the created_at and
// updated_at fields.
func (p *PackageAdd) MarshalJSON() ([]byte, error) {
	type Alias PackageAdd

	return json.Marshal(&struct {
		*Alias
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
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
	label := generateLabelFromName(name)

	return NewPackage(name, label)
}
