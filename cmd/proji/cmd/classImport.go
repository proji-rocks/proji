package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var directories, configs, exclude []string

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 {
			return fmt.Errorf("No flag was passed. You have to pass the '--config' or '--directory' flag at least once")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Import configs
		// Concat the two arrays so that '... import --config *.toml' is a valid command.
		// Without appending the args, proji would only use the first toml-file and not all of
		// them as intended with the '*'.
		// TODO: This section should be optimized and cleaned up.
		for _, config := range append(configs, args...) {
			if helper.IsInSlice(exclude, config) {
				continue
			}
			if err := importClassFromConfig(config, projiEnv.Svc); err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", config, err)
				continue
			}
			fmt.Printf("> '%s' was successfully imported\n", config)
		}

		// Import directories
		for _, directory := range directories {
			confName, err := importClassFromDirectory(directory, exclude, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", directory, err)
				continue
			}
			fmt.Printf("> Directory '%s' was successfully exported to '%s'\n", directory, confName)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)

	classImportCmd.Flags().StringSliceVar(&directories, "directory", []string{}, "import/imitate an existing directory")
	classImportCmd.MarkFlagDirname("directory")

	classImportCmd.Flags().StringSliceVar(&configs, "config", []string{}, "import a class from a config file")
	classImportCmd.MarkFlagFilename("config")

	classImportCmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "files/folders to exclude from import")
	classImportCmd.MarkFlagFilename("exclude")
}

func importClassFromConfig(config string, svc storage.Service) error {
	// Import class data
	class := item.NewClass("", "", false)
	if err := class.ImportFromConfig(config); err != nil {
		return err
	}
	return svc.SaveClass(class)
}

func importClassFromDirectory(directory string, excludeDir []string, svc storage.Service) (string, error) {
	// Import class data
	class := item.NewClass("", "", false)
	if err := class.ImportFromDirectory(directory, excludeDir); err != nil {
		return "", err
	}
	return class.Export()
}
