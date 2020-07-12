package domain

import (
	"time"

	"gorm.io/gorm"
)

// Template represents a template file or folder used by proji. It holds tags for gorm and toml defining its storage
// and export/import behaviour.
type Template struct {
	ID          uint           `gorm:"primarykey" toml:"-"`
	CreatedAt   time.Time      `toml:"-"`
	UpdatedAt   time.Time      `toml:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" toml:"-"`
	IsFile      bool           `gorm:"not null" toml:"is_file"`
	Path        string         `gorm:"index:idx_template_path_destination,unique;not null" toml:"path"`
	Destination string         `gorm:"index:idx_template_path_destination,unique;not null" toml:"destination"`
	Description string         `gorm:"size:255" toml:"description"`
}

func NewTemplate(path, destination, description string, isFile bool) *Template {
	return &Template{
		Path:        path,
		Destination: destination,
		Description: description,
		IsFile:      isFile,
	}
}
