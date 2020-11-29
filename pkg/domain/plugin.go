package domain

import (
	"time"
)

// Plugin represents a proji plugin that is may be used during the project creation process. It holds tags for gorm
// and toml defining its storage and export/import behaviour.
type Plugin struct {
	ID          uint      `gorm:"primarykey" toml:"-" json:"-"`
	CreatedAt   time.Time `toml:"-" json:"-"`
	UpdatedAt   time.Time `toml:"-" json:"-"`
	Path        string    `gorm:"index:idx_plugin_path,unique;not null" toml:"path" json:"path"`
	ExecNumber  int       `gorm:"check:(exec_number != 0);not null;size:4" toml:"exec_number" json:"exec_number"`
	Description string    `gorm:"size:255" toml:"description" json:"description"`
}
