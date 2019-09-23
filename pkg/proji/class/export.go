package class

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/viper"
)

// Export exports a given class to a toml config file
func (c *Class) Export() error {
	// Load the class data
	if err := c.Load(); err != nil {
		return err
	}

	// Create config string
	var configTxt = map[string]interface{}{
		"name":    c.Name,
		"labels":  c.Labels,
		"folders": c.Folders,
		"files":   c.Files,
		"scripts": c.Scripts,
	}

	// Export data to toml
	confName := "proji-export-" + c.Name + ".toml"
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

	destination, err := os.Create(destFolder + "/proji-class-example.toml")
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
