package storage

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Class struct represents a proji class
type Class struct {
	ID      uint              // Class ID in storage
	Name    string            // Class name
	Label   string            // Class label
	Folders map[string]string // Class related folders
	Files   map[string]string // Class related files
	Scripts map[string]bool   // Class related scripts
}

// NewClass returns a new class
func NewClass(name, label string) (*Class, error) {
	return &Class{
		ID:      0,
		Name:    name,
		Label:   label,
		Folders: make(map[string]string),
		Files:   make(map[string]string),
		Scripts: make(map[string]bool),
	}, nil
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
		"name":    c.Name,
		"label":   c.Label,
		"folders": c.Folders,
		"files":   c.Files,
		"scripts": c.Scripts,
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
