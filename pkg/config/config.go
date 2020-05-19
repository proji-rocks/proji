package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/nikoksr/proji/pkg/helper"
)

type configFile struct {
	src string
	dst string
}

type configFolder struct {
	path       string
	configs    []*configFile
	subFolders []string
}

// InitConfig is the main function for projis config initialization. It determines the OS' preferred config location, creates
// proji's config folders and downloads the required configs from GitHub to the local config folder.
func InitConfig(path, version string) (string, error) {
	fallbackVersion := "0.18.1"

	// Representation of default config folder
	cf := &configFolder{
		path: "",
		configs: []*configFile{
			{
				src: "https://raw.githubusercontent.com/nikoksr/proji/v" + version + "/assets/examples/example-config.toml",
				dst: "config.toml",
			},
			{
				src: "https://raw.githubusercontent.com/nikoksr/proji/v" + version + "/assets/examples/example-class-export.toml",
				dst: "examples/proji-class.toml",
			},
		},
		subFolders: []string{"db", "examples", "scripts", "templates"},
	}

	// Set OS specific config folder path
	cf.path = path

	// Create basefolder if it does not exist.
	err := helper.CreateFolderIfNotExists(cf.path)
	if err != nil {
		return "", err
	}

	// Create subfolders if they do not exist.
	for _, subFolder := range cf.subFolders {
		err = helper.CreateFolderIfNotExists(filepath.Join(cf.path, "/", subFolder))
		if err != nil {
			return "", err
		}
	}

	// Create configs if they do not exist.
	var wg sync.WaitGroup
	errors := make(chan error, len(cf.configs))

	for _, conf := range cf.configs {
		wg.Add(1)
		go func(conf *configFile, wg *sync.WaitGroup, e chan error) {
			defer wg.Done()
			dst := filepath.Join(cf.path, conf.dst)
			e <- helper.DownloadFileIfNotExists(conf.src, dst)
		}(conf, &wg, errors)
	}

	wg.Wait()
	close(errors)

	for err = range errors {
		if err != nil {
			if version == fallbackVersion {
				return cf.path, err
			}
			// Try with fallback version. This may help regular users but is manly for circleCI, which
			// fails when new versions are pushed. When a new version is pushed the corresponding github tag
			// doesn't exist, proji init fails.
			return InitConfig(cf.path, fallbackVersion)
		}
	}

	return cf.path, nil
}

// GetBaseConfigPath returns the OS specific path of the config folder.
func GetBaseConfigPath() (string, error) {
	configPath := ""
	switch runtime.GOOS {
	case "linux":
		home := os.Getenv("HOME")
		configPath = filepath.Join(home, "/.config/proji")
	case "darwin":
		home := os.Getenv("HOME")
		configPath = filepath.Join(home, "/Library/Application Support/proji")
	case "windows":
		appData := os.Getenv("APPDATA")
		configPath = filepath.Join(appData, "/proji")
	default:
		return "", fmt.Errorf("OS %s is not supported and/or tested yet. Please create an issue at "+
			"https://github.com/nikoksr/proji to request the support of your OS.\n", runtime.GOOS)
	}
	return configPath, nil
}
