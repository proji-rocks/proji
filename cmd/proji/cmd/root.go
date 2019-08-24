package cmd

import (
	"fmt"
	"os"

	"github.com/nikoksr/proji/internal/app/helper"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proji",
	Short: "Proji is a fast and simple project creator.",
	Long:  `Based on your favour proji creates the directory structures for a multitude of project types. Proji saves you hundrets of repetetive clicks or cli instructions. With one command proji creates you a multitude of projects based on your personal templates.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(helper.GetConfigDir())
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error: Could not read config file %s", viper.ConfigFileUsed())
		os.Exit(1)
	}
}
