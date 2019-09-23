package class

import (
	"github.com/BurntSushi/toml"
)

// Import imports a new class from a given config file.
func Import(configName string) (*Class, error) {
	var c Class
	if _, err := toml.DecodeFile(configName, &c); err != nil {
		return nil, err
	}

	return &c, c.Save()
}
