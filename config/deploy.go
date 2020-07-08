package config

import (
	"path/filepath"
	"sync"

	"github.com/nikoksr/proji/util"
)

// file represents a file that lives in the main config folder.
type file struct {
	src string
	dst string
}

// mainConfigFolder represents the structure of the main config folder.
type mainConfigFolder struct {
	basePath   string
	files      []*file
	subFolders []string
}

const (
	rawURLPrefix = "https://raw.githubusercontent.com/nikoksr/proji/v"
)

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

	// Write main config
	err = writeMainConfigFromDefaults(defaultConfigFolder.basePath)
	if err != nil {
		return err
	}

	// Download config files
	err = defaultConfigFolder.downloadFiles(
		version,
		fallbackVersion,
		forceUpdate,
	)
	return err
}

// newMainConfigFolder returns a struct that represents the structure of the main config folder.
func newMainConfigFolder() *mainConfigFolder {
	return &mainConfigFolder{
		basePath: "",
		files: []*file{
			{
				src: "/assets/examples/example-package-export.toml",
				dst: "examples/proji-package.toml",
			},
		},
		subFolders: []string{
			"db",
			"examples",
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

// downloadFiles downloads all files from github to the main config folder.
func (mcf *mainConfigFolder) downloadFiles(version, fallbackVersion string, forceUpdate bool) error {
	var wg sync.WaitGroup
	numFiles := len(mcf.files)
	wg.Add(numFiles)
	errs := make(chan error, numFiles)

	for _, conf := range mcf.files {
		go func(f *file) {
			defer wg.Done()
			src := rawURLPrefix + version + f.src
			dst := filepath.Join(mcf.basePath, f.dst)
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
				return mcf.downloadFiles(fallbackVersion, fallbackVersion, true)
			}
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
