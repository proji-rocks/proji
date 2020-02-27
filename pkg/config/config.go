package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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
				dst: "examples/class-export.toml",
			},
		},
		subFolders: []string{"db", "examples", "scripts", "templates"},
	}

	// Set OS specific config folder path
	cf.path = path

	// Create basefolder if it does not exist.
	err := createFolderIfNotExists(cf.path)
	if err != nil {
		return "", err
	}

	// Create subfolders if they do not exist.
	for _, subFolder := range cf.subFolders {
		err = createFolderIfNotExists(filepath.Join(cf.path, "/", subFolder))
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
			e <- downloadFileIfNotExists(conf.src, dst)
		}(conf, &wg, errors)
	}

	wg.Wait()
	close(errors)

	for err = range errors {
		if err != nil {
			if version == fallbackVersion {
				return cf.path, err
			} else {
				// Try with fallback version. This may help regular users but is manly for circleCI, which
				// fails when new versions are pushed. When a new version is pushed the corresponding github tag
				// doesn't exist, proji init fails.
				return InitConfig(cf.path, fallbackVersion)
			}
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

// createFolderIfNotExists creates a folder at the given path if it doesn't already exist.
func createFolderIfNotExists(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, os.ModePerm)
}

// downloadFile downloads a file from an url to the local fs.
func downloadFile(src, dst string) error {
	// Get the data
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error: %s", resp.Status)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// downloadFileIfNotExists runs downloadFile() if the destination file doesn't already exist.
func downloadFileIfNotExists(src, dst string) error {
	_, err := os.Stat(dst)
	if os.IsNotExist(err) {
		err = downloadFile(src, dst)
	}
	return err
}
