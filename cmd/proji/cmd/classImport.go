package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var directories, configs []string

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 {
			return fmt.Errorf("No config or directory path was given")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Import configs
		for _, config := range configs {
			if err := importClassFromConfig(config, projiEnv.Svc); err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", config, err)
				continue
			}
			fmt.Printf("> '%s' was successfully imported\n", config)
		}

		// Import directories
		for _, config := range directories {
			if err := importClassFromConfig(config, projiEnv.Svc); err != nil {
				fmt.Printf("> Import of '%s' failed: %v\n", config, err)
				continue
			}
			fmt.Printf("> '%s' was successfully imported\n", config)
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
}

func importClassFromConfig(config string, svc storage.Service) error {
	// Import class data
	class := item.NewClass("", "", false)
	if err := class.ImportFromConfig(config); err != nil {
		return err
	}
	return svc.SaveClass(class)
}

func importClassFromDirectory(directory string, svc storage.Service) error {
	return nil
}
