package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
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
func InitConfig(path, version string, forceUpdate bool) (string, error) {
	var err error

	if strings.Trim(path, " ") == "" {
		path, err = GetBaseConfigPath()
		if err != nil {
			return "", err
		}
	}

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
		go func(conf *configFile, forceUpdate bool, wg *sync.WaitGroup, e chan error) {
			defer wg.Done()
			dst := filepath.Join(cf.path, conf.dst)
			if forceUpdate {
				e <- helper.DownloadFile(conf.src, dst)
			} else {
				e <- helper.DownloadFileIfNotExists(conf.src, dst)
			}
		}(conf, forceUpdate, &wg, errors)
	}

	wg.Wait()
	close(errors)

	for err = range errors {
		if err != nil {
			if version == fallbackVersion {
				return cf.path, err
			}
			// Try with fallback version. This may help regular users but is manly for CI, which
			// fails when new versions are pushed. When a new version is pushed the corresponding github tag
			// doesn't exist, proji init fails.
			return InitConfig(cf.path, fallbackVersion, true)
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

func IsConfigUpToDate(projiVersion, configVersion string) (bool, error) {
	projiV, err := semver.NewVersion(projiVersion)
	if err != nil {
		return false, err
	}
	configV, err := semver.NewVersion(configVersion)
	if err != nil {
		return false, err
	}

	if configV.LessThan(projiV) {
		return false, errors.New("main config version is lower than proji's version. Please update your main" +
			" config")
	} else if configV.GreaterThan(projiV) {
		return true, errors.New("main config version is greater than proji's version, which could lead to " +
			"unforeseen errors")
	} else {
		return true, nil
	}
}

func ParsePathFromConfig(configFolderPath, pathToParse string) string {
	if filepath.IsAbs(pathToParse) {
		// Either user defined path like '/my/custom/db/path' or default value was loaded
		return pathToParse
	}
	// User defined path like 'db/proji.sqlite3'. Gets prefixed with config folder path. This has to be a relative
	// path or else the above will trigger.
	return filepath.Join(configFolderPath, pathToParse)
}
