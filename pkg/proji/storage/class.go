package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/viper"
)

// Class struct represents a proji class
type Class struct {
	// The class name
	Name string

	// The class ID
	ID uint

	// All class related labels
	Labels []string

	// All class related folders
	Folders map[string]string

	// All class related files
	Files map[string]string

	// All class related scripts
	Scripts map[string]bool
}

// NewClass returns a new class
func NewClass(name string) *Class {
	return &Class{
		Name:    name,
		ID:      0,
		Labels:  make([]string, 0),
		Folders: make(map[string]string),
		Files:   make(map[string]string),
		Scripts: make(map[string]bool),
	}
}

// Remove removes an existing class and all of its depending settings in other tables from the database.
func (c *Class) Remove(store Service) error {
	return store.RemoveClass(c.Name)
}

// ImportData imports class data from a given config file.
func (c *Class) ImportData(configName string) error {
	if _, err := toml.DecodeFile(configName, &c); err != nil {
		return err
	}
	return nil
}

// Export exports a given class to a toml config file
func (c *Class) Export() error {
	// Create config string
	var configTxt = map[string]interface{}{
		"name":    c.Name,
		"labels":  c.Labels,
		"folders": c.Folders,
		"files":   c.Files,
		"scripts": c.Scripts,
	}

	// Export data to toml
	confName := "proji-" + c.Name + ".toml"
	conf, err := os.Create(confName)
	if err != nil {
		return err
	}
	defer conf.Close()
	return toml.NewEncoder(conf).Encode(configTxt)
}

// ExportExample exports an example class config
func ExportExample(destFolder string) error {

	exampleDir, ok := viper.Get("examples.location").(string)
	if !ok {
		return fmt.Errorf("could not read example file location from config file")
	}
	exampleFile, ok := viper.Get("examples.class").(string)
	if !ok {
		return fmt.Errorf("could not read example file name from config file")
	}

	exampleFile = helper.GetConfigDir() + exampleDir + exampleFile
	sourceFileStat, err := os.Stat(exampleFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", exampleFile)
	}

	source, err := os.Open(exampleFile)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destFolder + "/proji-class.toml")
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
