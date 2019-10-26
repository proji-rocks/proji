package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exampleDest string
var exportAll bool

var classExportCmd = &cobra.Command{
	Use:   "export LABEL [LABEL...]",
	Short: "Export one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(exampleDest) > 0 {
			return exportExample(exampleDest, projiEnv.ConfPath)
		}

		if exportAll {
			err := exportAllClasses(projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Export of all classes failed: %v\n", err)
				return err
			}
			fmt.Println("> All classes were successfully exported")
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("Missing class label")
		}

		for _, label := range args {
			file, err := exportClass(label, projiEnv.Svc)
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
	classExportCmd.Flags().StringVarP(&exampleDest, "example", "e", "", "Export an example")
	classExportCmd.Flags().BoolVarP(&exportAll, "all", "a", false, "Export all classes")
}

func exportClass(label string, svc storage.Service) (string, error) {
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
	return class.Export()
}

func exportAllClasses(svc storage.Service) error {
	classes, err := svc.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		_, err = class.Export()
		if err != nil {
			return err
		}
	}
	return nil
}

func exportExample(destFolder, confPath string) error {
	examplePath, ok := viper.Get("examples.path").(string)
	if !ok {
		return fmt.Errorf("Could not read path of example config file")
	}

	examplePath = confPath + examplePath
	sourceFileStat, err := os.Stat(examplePath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", examplePath)
	}

	source, err := os.Open(examplePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(destFolder + "/proji-class-example.toml")
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
