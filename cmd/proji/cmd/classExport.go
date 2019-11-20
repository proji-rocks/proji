package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exportAll, example bool
var destination string

var classExportCmd = &cobra.Command{
	Use:   "export LABEL [LABEL...]",
	Short: "Export one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if exportAll && example {
			return fmt.Errorf("The flags 'example' and 'all' cannot be passed at the same time")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Export an example class
		if example {
			file, err := exportExample(destination, projiEnv.ConfPath)
			if err != nil {
				fmt.Printf("> Export of example class failed: %v\n", err)
				return err
			}
			fmt.Printf("> Example class was successfully exported to %s\n", file)
			return nil
		}

		// Export all classes
		if exportAll {
			err := exportAllClasses(destination, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Export of all classes failed: %v\n", err)
				return err
			}
			fmt.Println("> All classes were successfully exported")
			return nil
		}

		// Regular export
		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}

		for _, label := range args {
			file, err := exportClass(label, destination, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Export of '%s' to file %s failed: %v\n", label, file, err)
				continue
			}
			fmt.Printf("> '%s' was successfully exported to file %s\n", label, file)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classExportCmd)
	classExportCmd.Flags().BoolVarP(&example, "example", "e", false, "Export an example class")
	classExportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all classes")

	classExportCmd.Flags().StringVarP(&destination, "destination", "d", ".", "Destination for the export")
	classExportCmd.MarkFlagDirname("destination")
}

func exportClass(label, destination string, svc storage.Service) (string, error) {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return "", err
	}
	class, err := svc.LoadClass(classID)
	if err != nil {
		return "", err
	}
	if class.IsDefault {
		return "", fmt.Errorf("Default classes can not be exported")
	}
	return class.Export(destination)
}

func exportAllClasses(destination string, svc storage.Service) error {
	classes, err := svc.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		_, err = class.Export(destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func exportExample(destination, confPath string) (string, error) {
	examplePath, ok := viper.Get("examples.path").(string)
	if !ok {
		return "", fmt.Errorf("Could not read path of example config file")
	}

	examplePath = confPath + examplePath
	sourceFileStat, err := os.Stat(examplePath)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", examplePath)
	}

	src, err := os.Open(examplePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := destination + "/proji-class-example.toml"
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return dstPath, err
}
