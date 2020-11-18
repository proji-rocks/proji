package domain

import "time"

// Project represents a project that was created by proji. It holds tags for gorm and toml defining its storage and
// export/import behaviour.
type Project struct {
	ID        uint      `gorm:"primarykey" toml:"-"`
	CreatedAt time.Time `toml:"-"`
	UpdatedAt time.Time `toml:"-"`
	Name      string    `gorm:"size:64" toml:"name"`
	Path      string    `gorm:"index:idx_unq_project_path,unique;not null" toml:"path"`
	PackageID int       `toml:"-"`
	Package   *Package  `toml:"package"`
}

func NewProject(name, path string, pkg *Package) *Project {
	return &Project{
		Name:    name,
		Path:    path,
		Package: pkg,
	}
}

type ProjectStore interface {
	StoreProject(p *Project) error

	LoadProject(path string) (*Project, error)
	LoadProjectList(paths ...string) ([]*Project, error)

	UpdateProjectLocation(oldPath, newPath string) error

	RemoveProject(path string) error
}

type ProjectService interface {
	StoreProject(p *Project) error
	LoadProject(path string) (*Project, error)
	LoadProjectList(paths ...string) ([]*Project, error)
	UpdateProjectLocation(oldPath, newPath string) error
	RemoveProject(path string) error

	CreateProject(configRootPath string, project *Project) (err error)
}
