package domain

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project that was created by proji. It holds tags for gorm and toml defining its storage and
// export/import behaviour.
type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_unq_project_path_deletedat,unique;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:64"`
	Path      string         `gorm:"index:idx_unq_project_path_deletedat,unique;not null"`
	Package   *Package       `gorm:"ForeignKey:ID;References:ID"`
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
	PurgeProject(path string) error
}

type ProjectService interface {
	StoreProject(p *Project) error
	LoadProject(path string) (*Project, error)
	LoadProjectList(paths ...string) ([]*Project, error)
	UpdateProjectLocation(oldPath, newPath string) error
	RemoveProject(path string) error
	PurgeProject(path string) error

	CreateProject(configRootPath string, project *Project) (err error)
}
