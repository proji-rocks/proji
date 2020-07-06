package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/config"

	"github.com/nikoksr/proji/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Env represents central resources and information the app uses.
type Env struct {
	Auth             *config.APIAuthentication
	DBPath           string
	Svc              storage.Service
	ConfigFolderPath string
	ExcludedPaths    []string
	FallbackVersion  string
	Version          string
}

var projiEnv *Env

const (
	configExcludeFoldersKey = "import.exclude_folders"
	configDBKey             = "sqlite3.path"
	configGHTokenKey        = "auth.gh_token" //nolint:gosec
	configGLTokenKey        = "auth.gl_token" //nolint:gosec
)

var rootCmd = &cobra.Command{
	Use:   "proji",
	Short: "A powerful cross-platform CLI project templating tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if projiEnv == nil {
			log.Fatalf("Error: Env struct is not defined.\n")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	if projiEnv == nil {
		projiEnv = &Env{
			Auth:             &config.APIAuthentication{},
			DBPath:           "",
			ExcludedPaths:    make([]string, 0),
			ConfigFolderPath: "",
			Svc:              nil,
			FallbackVersion:  "0.19.2",
			Version:          "0.20.0",
		}
	}

	if len(os.Args) > 1 && os.Args[1] != "init" && os.Args[1] != "version" && os.Args[1] != "help" {
		cobra.OnInitialize(initConfig, initStorageService)
	}
}

func initConfig() {
	// Set platform specific config path
	var err error
	projiEnv.ConfigFolderPath, err = config.GetBaseConfigPath()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	viper.AddConfigPath(projiEnv.ConfigFolderPath)

	// Config name
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// Set default values as fallback
	setDefaultConfigValues()

	// Read config
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error: %v\n\nTry and execute: proji init\n", err)
	}

	// Read environment variables
	viper.SetEnvPrefix("PROJI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set all proji relevant environmental values
	setAllEnvValues()
}

func initStorageService() {
	var err error
	/*
		projiEnv.Svc, err = sqlite.New(projiEnv.DBPath)
		if err != nil {
			log.Fatalf("Error: could not connect to sqlite db. %v\n%s\n", err, projiEnv.DBPath)
		}
	*/
	projiEnv.Svc, err = storage.NewService("sqlite3", "/home/niko/tmp.db")
	if err != nil {
		log.Fatalf("Error: could not connect to sqlite db. %v\n%s\n", err, projiEnv.DBPath)
	}
}

func setDefaultConfigValues() {
	viper.SetDefault(configGHTokenKey, "")
	viper.SetDefault(configGLTokenKey, "")
	viper.SetDefault(configExcludeFoldersKey, make([]string, 0))
	viper.SetDefault(configDBKey, filepath.Join(projiEnv.ConfigFolderPath, "/db/proji.sqlite3"))
}

func setAllEnvValues() {
	projiEnv.Auth.GHToken = viper.GetString(configGHTokenKey)
	projiEnv.Auth.GLToken = viper.GetString(configGLTokenKey)
	projiEnv.ExcludedPaths = viper.GetStringSlice(configExcludeFoldersKey)
	projiEnv.DBPath = config.ParsePathFromConfig(projiEnv.ConfigFolderPath, viper.GetString(configDBKey))
}
