package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Env represents central resources and information the app uses.
type env struct {
	Svc      storage.Service
	ConfPath string
}

var projiEnv *env

var rootCmd = &cobra.Command{
	Use:   "proji",
	Short: "A lean and mean project creator and manager.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		svc, err := initStorageService()
		if err != nil {
			return err
		}
		if projiEnv == nil {
			return fmt.Errorf("env struct is not defined")
		}
		projiEnv.Svc = svc
		return nil
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
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Could not set config path: %v", err)
		os.Exit(1)
	}
	if projiEnv == nil {
		projiEnv = &env{ConfPath: "", Svc: nil}
	}
	projiEnv.ConfPath = home + "/.config/proji/"
	viper.AddConfigPath(projiEnv.ConfPath)
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Could not read config file %s", viper.ConfigFileUsed())
		os.Exit(1)
	}
}

func initStorageService() (storage.Service, error) {
	dbPath, ok := viper.Get("sqlite3.path").(string)
	if !ok {
		return nil, fmt.Errorf("could not read sqlite path from config file")
	}

	svc, err := sqlite.New(projiEnv.ConfPath + dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not connect to sqlite db: %v", err)
	}
	return svc, nil
}
