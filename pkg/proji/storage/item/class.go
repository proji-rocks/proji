package item

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Class struct represents a proji class
type Class struct {
	ID        uint      // Class ID in storage. Not exported/imported.
	IsDefault bool      // Is this a default class? Not exported/imported.
	Name      string    `toml:"name"`   // Class name
	Label     string    `toml:"label"`  // Class label
	Folders   []*Folder `toml:"folder"` // Class folders
	Files     []*File   `toml:"file"`   // Class files
	Scripts   []*Script `toml:"script"` // Class scripts
}

// NewClass returns a new class
func NewClass(name, label string, isDefault bool) *Class {
	return &Class{
		ID:        0,
		IsDefault: isDefault,
		Name:      name,
		Label:     label,
		Folders:   make([]*Folder, 0),
		Files:     make([]*File, 0),
		Scripts:   make([]*Script, 0),
	}
}

// ImportData imports class data from a given config file.
func (c *Class) ImportData(configName string) error {
	// Validate that it's a toml file
	if !strings.HasSuffix(configName, ".toml") {
		return fmt.Errorf("Import file has to be of type 'toml'")
	}

	// Validate config is not empty
	conf, err := os.Stat(configName)
	if err != nil {
		return err
	}
	if conf.Size() == 0 {
		return fmt.Errorf("Import file is empty")
	}

	// Decode the file
	_, err = toml.DecodeFile(configName, &c)
	return err
}

// Export exports a given class to a toml config file
func (c *Class) Export() (string, error) {
	// Create config string
	var configTxt = map[string]interface{}{
		"name":   c.Name,
		"label":  c.Label,
		"folder": c.Folders,
		"file":   c.Files,
		"script": c.Scripts,
	}

	// Export data to toml
	confName := "proji-" + c.Name + ".toml"
	conf, err := os.Create(confName)
	if err != nil {
		return confName, err
	}
	defer conf.Close()
	return confName, toml.NewEncoder(conf).Encode(configTxt)
}
