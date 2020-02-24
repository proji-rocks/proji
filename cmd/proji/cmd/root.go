package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/config"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Env represents central resources and information the app uses.
type env struct {
	Svc            storage.Service
	UserConfigPath string
	Excludes       []string
	Version        string
}

var projiEnv *env

var rootCmd = &cobra.Command{
	Use:   "proji",
	Short: "A fast and powerful cli project scaffolding tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if projiEnv == nil {
			log.Fatalf("Error: env struct is not defined.\n")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if projiEnv.Svc != nil {
			_ = projiEnv.Svc.Close()
			projiEnv.Svc = nil
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
		projiEnv = &env{Svc: nil, UserConfigPath: "", Excludes: make([]string, 0), Version: "0.18.1"}
	}

	var err error
	projiEnv.UserConfigPath, err = config.GetBaseConfigPath()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	viper.AddConfigPath(projiEnv.UserConfigPath)
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if len(os.Args) > 1 && os.Args[1] != "init" && os.Args[1] != "version" && os.Args[1] != "help" {
		cobra.OnInitialize(initConfig, initStorageService)
	}
}

func initConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error: %v\n\nTry and execute: proji init\n", err)
	}
	projiEnv.Excludes = viper.GetStringSlice("import.excludeFolders")
}

func initStorageService() {
	dbPath := viper.GetString("sqlite3.path")
	var err error
	projiEnv.Svc, err = sqlite.New(filepath.Join(projiEnv.UserConfigPath, dbPath))
	if err != nil {
		log.Fatalf("Error: could not connect to sqlite db. %v\n%s\n", err, projiEnv.UserConfigPath)
	}
}
