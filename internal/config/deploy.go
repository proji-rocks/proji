package config

import (
	"path/filepath"

	"github.com/nikoksr/proji/internal/util"
)

// mainConfigFolder represents the structure of the main config folder.
type mainConfigFolder struct {
	basePath   string
	subFolders []string
}

// Deploy deploys projis main config folder to disk. It creates all subfolders, downloads needed files and
// creates a main config file from default values. Version determines from which version of proji
// files should be downloaded. If version fails, it will try again and use the fallback version.
// ForceUpdate should usually not be used and is only used internally to overwrite existing files if necessary.
func Deploy(version, fallbackVersion string, forceUpdate bool) error {
	defaultConfigFolder := newMainConfigFolder()
	defaultConfigFolder.basePath = globalBasePath

	// Create subfolders
	err := defaultConfigFolder.createSubFolders()
	if err != nil {
		return err
	}

	// Write main config from defaults
	return writeMainConfigFromDefaults(defaultConfigFolder.basePath)
}

// newMainConfigFolder returns a struct that represents the structure of the main config folder.
func newMainConfigFolder() *mainConfigFolder {
	return &mainConfigFolder{
		basePath: "",
		subFolders: []string{
			"db",
			"plugins",
			"templates",
		},
	}
}

// createSubFolders creates a list of subfolders in the main config folder.
func (mcf *mainConfigFolder) createSubFolders() error {
	for _, subFolder := range mcf.subFolders {
		err := util.CreateFolderIfNotExists(filepath.Join(mcf.basePath, subFolder))
		if err != nil {
			return err
		}
	}
	return nil
}

// writeMainConfigFromDefaults creates a config instance with default values and writes it to the given path.
func writeMainConfigFromDefaults(path string) error {
	conf := New(path)
	conf.setProvider()
	conf.setSpecs()
	conf.setDefaultValues()
	return conf.provider.SafeWriteConfig()
}
