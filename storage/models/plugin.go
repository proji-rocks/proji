package models

import (
	"time"

	"gorm.io/gorm"
)

// Plugin represents a proji plugin that is may be used during the project creation process. It holds tags for gorm
// and toml defining its storage and export/import behaviour.
type Plugin struct {
	ID          uint           `gorm:"primarykey" toml:"-"`
	CreatedAt   time.Time      `toml:"-"`
	UpdatedAt   time.Time      `toml:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" toml:"-"`
	Path        string         `gorm:"index:idx_plugin_path,unique;not null" toml:"path"`
	ExecNumber  int            `gorm:"check:(exec_number != 0);not null;size:4" toml:"exec_number"`
	Description string         `gorm:"size:255" toml:"description"`
}
