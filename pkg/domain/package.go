package domain

import (
	"net/url"
	"regexp"
	"time"
)

// Package represents a proji package; the central item of proji's project creation mechanism. It holds tags for gorm and
// toml defining its storage and export/import behaviour.
type Package struct {
	ID          uint        `gorm:"primarykey" toml:"-" json:"-"`
	CreatedAt   time.Time   `toml:"-" json:"-"`
	UpdatedAt   time.Time   `toml:"-" json:"-"`
	Name        string      `gorm:"not null;size:64" toml:"name" json:"name"`
	Label       string      `gorm:"index:idx_unq_package_label,unique;not null;size:16" toml:"label" json:"label"`
	Description string      `gorm:"size:255" toml:"description" json:"description"`
	Templates   []*Template `gorm:"many2many:package_templates;" toml:"template" json:"template"`
	Plugins     []*Plugin   `gorm:"many2many:package_plugins;" toml:"plugin" json:"plugin"`
}

func NewPackage(name, label string) *Package {
	return &Package{
		Name:  name,
		Label: label,
	}
}

type PackageStore interface {
	StorePackage(p *Package) error

	LoadPackage(loadDependencies bool, label string) (*Package, error)
	LoadPackageList(loadDependencies bool, labels ...string) ([]*Package, error)

	RemovePackage(label string) error
}

type PackageService interface {
	StorePackage(p *Package) error
	LoadPackage(loadDependencies bool, label string) (*Package, error)
	LoadPackageList(loadDependencies bool, labels ...string) ([]*Package, error)
	RemovePackage(label string) error

	ImportPackageFromConfig(path string) (*Package, error)
	ImportPackageFromDirectoryStructure(path string, exclude *regexp.Regexp) (*Package, error)
	ImportPackageFromRepositoryStructure(url *url.URL, exclude *regexp.Regexp) (*Package, error)
	ImportPackageFromRemote(url *url.URL) (*Package, error)
	ImportPackagesFromCollection(url *url.URL, exclude *regexp.Regexp) ([]*Package, error)

	ExportPackageToConfig(pkg Package, destination string, asJson bool) (string, error)
}
