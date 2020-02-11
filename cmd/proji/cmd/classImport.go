package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/spf13/cobra"
)

var remoteRepos, directories, configs, excludes []string

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 && len(remoteRepos) < 1 {
			return fmt.Errorf("no flag was passed. You have to pass the '--config', '--directory' or '--remote-repo' flag at least once")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Concat the two arrays so that '... import --config *.toml' is a valid command.
		// Without appending the args, proji would only use the first toml-file and not all of
		// them as intended with the '*'.
		configs = append(configs, args...)
		pathMap := map[string][]string{"files": configs, "dirs": directories, "urls": remoteRepos}

		// Import configs
		for pathType, paths := range pathMap {
			for _, path := range paths {
				result, err := importClass(path, pathType, excludes)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Println(result)
				}
			}
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)

	classImportCmd.Flags().StringSliceVar(&remoteRepos, "remote-repo", []string{}, "create an importable config based on a remote repository")
	_ = classImportCmd.MarkFlagDirname("remote-repo")

	classImportCmd.Flags().StringSliceVar(&directories, "directory", []string{}, "create an importable config based on a local directory")
	_ = classImportCmd.MarkFlagDirname("directory")

	classImportCmd.Flags().StringSliceVar(&configs, "config", []string{}, "import a class from a config file")
	_ = classImportCmd.MarkFlagFilename("config")

	classImportCmd.Flags().StringSliceVar(&excludes, "exclude", []string{}, "folder to exclude from local directory import")
	_ = classImportCmd.MarkFlagFilename("exclude")
}

func importClass(path, pathType string, excludes []string) (string, error) {
	if helper.IsInSlice(excludes, path) {
		return "", nil
	}

	class := item.NewClass("", "", false)
	var err error
	var confName, msg string

	switch pathType {
	case "files":
		err = class.ImportFromConfig(path)
		if err != nil {
			return "", err
		}
		err = projiEnv.Svc.SaveClass(class)
		if err == nil {
			msg = fmt.Sprintf("> Successfully imported class '%s' from '%s'", class.Name, path)
		}
	case "dirs":
		fallthrough
	case "urls":
		if pathType == "dirs" {
			err = class.ImportFromDirectory(path, excludes)
			if err != nil {
				return "", err
			}
		} else {
			err = class.ImportFromURL(path)
			if err != nil {
				return "", err
			}
		}
		// Classes that are generated from directories and URLs should be exported to a config file first
		// so that the user can fine tune them
		confName, err = class.Export(".")
		if err == nil {
			msg = fmt.Sprintf("> '%s' was successfully exported to '%s'", path, confName)
		}
	default:
		err = fmt.Errorf("path type %s is not supported", pathType)
	}
	return msg, err
}
