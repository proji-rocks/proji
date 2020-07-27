package domain

import (
	"time"
)

// Template represents a template file or folder used by proji. It holds tags for gorm and toml defining its storage
// and export/import behaviour.
type Template struct {
	ID          uint      `gorm:"primarykey" toml:"-"`
	CreatedAt   time.Time `toml:"-"`
	UpdatedAt   time.Time `toml:"-"`
	IsFile      bool      `gorm:"not null" toml:"is_file"`
	Destination string    `gorm:"index:idx_template_path_destination,unique;not null" toml:"destination"`
	Path        string    `gorm:"index:idx_template_path_destination,unique;not null" toml:"path"`
	Description string    `gorm:"size:255" toml:"description"`
}
