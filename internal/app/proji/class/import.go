package class

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// classConfig represents a class config file
type classConfig struct {
	Title   string
	Class   map[string]string
	Labels  map[string][]string
	Folders map[string]string
	Files   map[string]string
	Scripts map[string]bool
}

// Import imports a new class from a given config file.
func Import(configName string) error {
	var conf classConfig
	if _, err := toml.DecodeFile(configName, &conf); err != nil {
		return err
	}

	fmt.Printf("> Importing %s...\n", conf.Title)
	err := AddClassToDB(conf.Class["name"], conf.Labels["data"], conf.Folders, conf.Files, conf.Scripts)
	if err != nil {
		return err
	}
	fmt.Println("> Done...")
	return nil
}
