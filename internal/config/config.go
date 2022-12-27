package config

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	// Auth is a general authentication configuration.
	Auth struct {
		// GitHubToken allows proji to access private GitHub repos and avoid rate limiting.
		GitHubToken string `mapstructure:"github_token"`
		// GitLabToken allows proji to access private GitLab repos and avoid rate limiting.
		GitLabToken string `mapstructure:"gitlab_token"`
	}

	// Database is a configuration for the database.
	Database struct {
		// DSN is the data source name for the database. At the moment, this is a path to a boltDB.
		DSN string `mapstructure:"dsn"`
	}

	// Import is a configuration for the import.
	Import struct {
		// Exclude is a regex that is used to exclude files and directories from the import process.
		Exclude string `mapstructure:"exclude"`
	}

	// Sentry is a configuration for the Sentry monitoring service.
	Sentry struct {
		// Enabled is a flag that indicates if Sentry is enabled. This is disabled by default.
		Enabled bool `mapstructure:"enabled"`
	}

	// Monitoring returns the monitoring configuration.
	Monitoring struct {
		// Sentry is a configuration for the Sentry monitoring service.
		Sentry Sentry `mapstructure:"sentry"`
	}

	// System includes settings that in some way or another affect the system that proji is running on.
	System struct {
		// TextEditor is the editor that proji will use to open files. This is usually set to the $EDITOR env var.
		TextEditor string `mapstructure:"text_editor"`
	}

	// Config is the configuration for the application.
	Config struct {
		Auth       Auth         `mapstructure:"-"`
		Database   Database     `mapstructure:"database"`
		Import     Import       `mapstructure:"import"`
		Monitoring Monitoring   `mapstructure:"monitoring"`
		System     System       `mapstructure:"system"`
		provider   *viper.Viper `mapstructure:"-"`
	}
)

const (
	// Config directory defaults
	defaultConfigDir  = "proji"
	defaultConfigFile = "config.toml"

	// Data directory defaults
	defaultDataDir      = "data"
	defaultDatabaseFile = "proji.db"

	// Other subdirectories
	defaultPluginsDir   = "plugins"
	defaultTemplatesDir = "templates"

	// Some config constants/defaults
	defaultExcludePattern = `^(.git|.env|.idea|.vscode)$`
	defaultSentryState    = false
)

var (
	// Anonymous check to ensure that the default regex exclude pattern is valid.
	_ = regexp.MustCompile(defaultExcludePattern)

	config   *Config   // Singleton
	loadOnce sync.Once // Used to ensure the singleton is in fact only loaded once

	// ErrUnsupportedOS is returned when the operating system is not supported.
	ErrUnsupportedOS = errors.New("unsupported operating system")

	// ErrInvalidUserConfigPath is returned when the user config path could not be determined.
	ErrInvalidUserConfigPath = errors.New("could not determine user config path")
)

// getUnixUserConfigPath returns the path to the user config directory for *nix systems. This uses the $HOME env var.
func getUnixUserConfigPath() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return "", ErrInvalidUserConfigPath
	}

	return filepath.Join(homeDir, ".config", defaultConfigDir, defaultConfigFile), nil
}

// getWindowsUserConfigPath returns the path to the user config directory for Windows systems. This uses the
// %LOCALAPPDATA% env var.
func getWindowsUserConfigPath() (string, error) {
	homeDir := os.Getenv("LOCALAPPDATA")
	if homeDir == "" {
		return "", ErrInvalidUserConfigPath
	}

	return filepath.Join(homeDir, defaultConfigDir, defaultConfigFile), nil
}

// getUserConfigPath returns the path to the user config directory for the given operating system. We pass in the
// operating system to make it easier to test.
func getUserConfigPath(system string) (path string, err error) {
	switch system {
	// Same path for *nix systems
	case "darwin", "freebsd", "linux", "netbsd", "openbsd":
		path, err = getUnixUserConfigPath()
	case "windows":
		path, err = getWindowsUserConfigPath()
	default:
		return "", ErrUnsupportedOS
	}

	return path, err
}

// defaultConfigPath returns the default path to the configuration file. This usually gets called when no configuration
// path is provided by the user.
func defaultConfigPath() (string, error) {
	return getUserConfigPath(runtime.GOOS)
}

// newProvider creates a new viper instance and sets the default values.
func newProvider(path string) *viper.Viper {
	provider := viper.New()

	// Allow for cross-platform paths
	path = filepath.Clean(path)
	dir := filepath.Dir(path)

	// Set default configuration
	provider.SetDefault("database.dsn", filepath.Join(dir, defaultDataDir, defaultDatabaseFile))
	provider.SetDefault("import.exclude", defaultExcludePattern)
	provider.SetDefault("monitoring.sentry.enabled", defaultSentryState)

	// Set configuration file path
	provider.SetConfigFile(path)

	return provider
}

// setupInfrastructure sets up the infrastructure for the configuration. It creates all necessary files and directories.
func (conf *Config) setupInfrastructure() error {
	if conf.provider == nil {
		return errors.New("config provider is nil")
	}

	// Get directory for config file
	configPath := conf.provider.ConfigFileUsed()
	baseDir := filepath.Dir(configPath)

	// Create subdirectories; this also implicitly creates the base directory
	dirs := []string{
		filepath.Join(baseDir, defaultPluginsDir),
		filepath.Join(baseDir, defaultTemplatesDir),
		filepath.Join(baseDir, defaultDataDir),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return errors.Wrapf(err, "create directory %q", dir)
		}
	}

	// Create config file if it does not exist
	err := conf.provider.SafeWriteConfigAs(configPath)
	if err != nil {
		return errors.Wrapf(err, "write config to %q", configPath)
	}

	// TODO: This is declared as a DSN, but it is actually a file path. I want to keep proji open to support more
	//       databases in the future, so DSN is probably right, but the way we handle it here is not. We should use
	//       a different approach. Check the protocol of the dsn or something similar. If it is not sqlite or bolt, we
	//       should not create the file.
	// Create empty database file
	dbFile := conf.provider.GetString("database.dsn")
	_, err = os.Create(dbFile)

	return errors.Wrapf(err, "create file %q", dbFile)
}

// readFile loads configuration from the configuration file.
func (conf *Config) readFile() error {
	return conf.provider.ReadInConfig()
}

// readEnv loads configuration from environment variables.
func (conf *Config) readEnv() {
	conf.provider.AutomaticEnv()
}

// readFlags loads configuration from command line flags.
func (conf *Config) readFlags(cmdFlags *pflag.FlagSet) error {
	if cmdFlags == nil {
		return nil
	}

	// Map flags to viper config tags
	flags := map[string]string{
		"exclude": "import.exclude",
	}

	for name, key := range flags {
		flag := cmdFlags.Lookup(name)
		if flag == nil {
			continue // Flag not set
		}

		err := conf.provider.BindPFlag(key, flag)
		if err != nil {
			return errors.Wrap(err, "bind flag")
		}
	}

	return nil
}

// load the configuration. It combines loading from the configuration file, environment variables and command line
// flags.
// If the path is empty, the default configuration path is used. If the path is not empty, the path must point directly
// to the configuration file.
func load(ctx context.Context, path string, flags *pflag.FlagSet) (conf *Config, err error) {
	logger := simplog.FromContext(ctx)

	// If no explicit path is given, use default path
	if path == "" {
		path, err = defaultConfigPath()
		if err != nil {
			return nil, errors.Wrap(err, "get config path")
		}

		logger.Debugf("no explicit config path given, using default path: %q", path)
	}

	// Clean up path
	path = filepath.Clean(path)
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "get absolute config path")
	}

	// Create default config
	logger.Debugf("creating config provider with path: %q", path)
	conf = &Config{
		provider: newProvider(path),
	}

	// Load config values; order is important here. File < Env < Flags.
	logger.Debugf("loading config values from file: %q", path)
	err = conf.readFile()

	// If the config file doesn't exist, create it.
	if errors.Is(err, os.ErrNotExist) {
		logger.Debugf("config file does not exist, setting up infrastructure : %q", filepath.Dir(path))
		err = conf.setupInfrastructure()
		if err != nil {
			return nil, errors.Wrap(err, "setup infrastructure")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	// Environment variables
	logger.Debugf("loading config values from environment")
	conf.readEnv()

	// Command line flags
	logger.Debugf("loading config values from flags")
	err = conf.readFlags(flags)
	if err != nil {
		return conf, errors.Wrap(err, "load config from flags")
	}

	// Unmarshal config values into struct
	logger.Debugf("unmarshalling config")
	err = conf.provider.Unmarshal(conf)
	if err != nil {
		return conf, errors.Wrap(err, "unmarshal config")
	}

	return conf, nil
}

// Load loads the configuration from the given path. If the path is empty, the default path is used. If a path is given,
// it has to be an absolute path that points directly to the config file. It returns the singleton instance.
func Load(ctx context.Context, path string, flags *pflag.FlagSet) (*Config, error) {
	var err error
	loadOnce.Do(func() {
		config, err = load(ctx, path, flags)
		if err != nil {
			return
		}
	})

	return config, err
}

// Validate the configuration. Invalid configs will return a ErrInvalidConfig error wrapped in a more detailed error.
// This method never gets called automatically. It has to be called manually.
func (conf *Config) Validate() error {
	// Validate database
	if conf.Database.DSN == "" {
		return errors.New("database dsn is empty")
	}

	return nil
}

// BaseDir returns the base directory of the configuration file.
func (conf *Config) BaseDir() string {
	return filepath.Dir(conf.provider.ConfigFileUsed())
}

// PluginsDir returns the plugins' directory.
func (conf *Config) PluginsDir() string {
	return filepath.Join(conf.BaseDir(), defaultPluginsDir)
}

// TemplatesDir returns the templates' directory.
func (conf *Config) TemplatesDir() string {
	return filepath.Join(conf.BaseDir(), defaultTemplatesDir)
}
