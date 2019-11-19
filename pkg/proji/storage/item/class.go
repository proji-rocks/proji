package item

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/nikoksr/proji/pkg/helper"

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

// ImportFromConfig imports class data from a given config file.
func (c *Class) ImportFromConfig(configName string) error {
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

	if len(c.Name) < 1 {
		return fmt.Errorf("Name cannot be an empty string")
	}
	if len(c.Label) < 1 {
		return fmt.Errorf("Label cannot be an empty string")
	}

	return err
}

// ImportFromDirectory imports a class from a given directory. Proji will copy the
// structure and content of the directory and create a class based on it.
func (c *Class) ImportFromDirectory(directory string, excludeDirs []string) error {
	// Validate that the directory exists
	if !helper.DoesPathExist(directory) {
		return fmt.Errorf("Given directory does not exist")
	}

	// Set class name from directory base name
	base := path.Base(directory)
	c.Name = base
	c.Label = pickLabel(c.Name)

	// This map of directories that should be skipped might be moved to the main config
	// file so that it's editable and extensible.
	excludeDirs = append(excludeDirs, []string{".git", ".env"}...)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// Skip base directory
		if directory == path {
			return nil
		}
		// Extract relative path
		relPath, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		// Add file or folder to class
		if info.IsDir() {
			c.Folders = append(c.Folders, &Folder{Destination: relPath, Template: ""})
			if helper.IsInSlice(excludeDirs, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			c.Files = append(c.Files, &File{Destination: relPath, Template: ""})
		}
		return nil
	})
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

func pickLabel(className string) string {
	nameLen := len(className)
	if nameLen < 2 {
		return strings.ToLower(className)
	}

	label := ""
	maxLabelLen := 4

	// Try to create label by separators
	seps := []string{"-", "_", ".", " "}
	parts := make([]string, 0)

	for _, d := range seps {
		parts = strings.Split(className, d)
		if len(parts) > 1 {
			break
		}
	}

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Try to create label by uppercase letters
	if !unicode.IsUpper(rune(className[0])) {
		className = string(byte(unicode.ToUpper(rune(className[0])))) + className[1:]
	}

	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	parts = re.FindAllString(className, -1)

	if len(parts) > 1 {
		for i, part := range parts {
			if i > maxLabelLen {
				break
			}
			label += string(part[0])
		}
		return strings.ToLower(label)
	}

	// Pick first, mid and last byte in string
	label = string(className[0]) + string(className[nameLen/2]) + string(className[nameLen-1])
	return strings.ToLower(label)
}
