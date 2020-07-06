package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/nikoksr/proji/util"
)

type APIAuthentication struct {
	GHToken string
	GLToken string
}

type configFile struct {
	src string
	dst string
}

type mainConfigFolder struct {
	basePath   string
	configs    []*configFile
	subFolders []string
}

const (
	rawURLPrefix = "https://raw.githubusercontent.com/nikoksr/proji/v"
)

// InitConfig is the main function for projis config initialization. It determines the OS' preferred config location, creates
// proji's config folders and downloads the required configs from GitHub to the local config folder.
func InitConfig(path, version, fallbackVersion string, forceUpdate bool) error {
	var err error

	// Set base config path if not given
	if strings.Trim(path, " ") == "" {
		path, err = GetBaseConfigPath()
		if err != nil {
			return err
		}
	}

	// Representation of proji's main config folder
	defaultConfigFolder := &mainConfigFolder{
		basePath: "",
		configs: []*configFile{
			{
				src: "/assets/examples/example-config.toml",
				dst: "config.toml",
			},
			{
				src: "/assets/examples/example-class-export.toml",
				dst: "examples/proji-class.toml",
			},
		},
		subFolders: []string{"db", "examples", "scripts", "templates"},
	}

	defaultConfigFolder.basePath = path

	// Create subfolders
	err = defaultConfigFolder.createSubFolders()
	if err != nil {
		return err
	}

	// Download config files
	err = defaultConfigFolder.downloadConfigFiles(
		version,
		fallbackVersion,
		forceUpdate,
	)
	return err
}

// Create subfolders if they do not exist.
func (mcf *mainConfigFolder) createSubFolders() error {
	for _, subFolder := range mcf.subFolders {
		err := util.CreateFolderIfNotExists(filepath.Join(mcf.basePath, subFolder))
		if err != nil {
			return err
		}
	}
	return nil
}

func (mcf *mainConfigFolder) downloadConfigFiles(version, fallbackVersion string, forceUpdate bool) error {
	var wg sync.WaitGroup
	numConfigs := len(mcf.configs)
	wg.Add(numConfigs)
	errs := make(chan error, numConfigs)

	for _, conf := range mcf.configs {
		go func(conf *configFile) {
			defer wg.Done()
			src := rawURLPrefix + version + conf.src
			dst := filepath.Join(mcf.basePath, conf.dst)
			if forceUpdate {
				errs <- util.DownloadFile(dst, src)
			} else {
				errs <- util.DownloadFileIfNotExists(dst, src)
			}
		}(conf)
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			if version != fallbackVersion {
				// Try with fallback version. This may help regular users but is manly for CI, which
				// fails when new versions are pushed. When a new version is pushed the corresponding github tag
				// doesn't exist, proji init fails.
				return mcf.downloadConfigFiles(fallbackVersion, fallbackVersion, true)
			}
			return err
		}
	}
	return nil
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
			"https://github.com/nikoksr/proji to request the support of your OS", runtime.GOOS)
	}
	return configPath, nil
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
