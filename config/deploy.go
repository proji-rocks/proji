package config

import (
	"path/filepath"
	"sync"

	"github.com/nikoksr/proji/util"
)

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
func Deploy(version, fallbackVersion string, forceUpdate bool) error {
	defaultConfigFolder := newMainConfigFolder()
	defaultConfigFolder.basePath = globalBasePath

	// Create subfolders
	err := defaultConfigFolder.createSubFolders()
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

func newMainConfigFolder() *mainConfigFolder {
	return &mainConfigFolder{
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
		subFolders: []string{"db", "examples", "plugins", "templates"},
	}
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
