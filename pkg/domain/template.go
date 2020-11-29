package domain

import (
	"time"
)

// Template represents a template file or folder used by proji. It holds tags for gorm and toml defining its storage
// and export/import behaviour.
type Template struct {
	ID          uint      `gorm:"primarykey" toml:"-" json:"-"`
	CreatedAt   time.Time `toml:"-" json:"-"`
	UpdatedAt   time.Time `toml:"-" json:"-"`
	IsFile      bool      `gorm:"not null" toml:"is_file" json:"is_file"`
	Destination string    `gorm:"index:idx_template_path_destination,unique;not null" toml:"destination" json:"destination"`
	Path        string    `gorm:"index:idx_template_path_destination,unique;not null" toml:"path" json:"path"`
	Description string    `gorm:"size:255" toml:"description" json:"description"`
}
