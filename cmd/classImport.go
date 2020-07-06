package cmd

import (
	"fmt"
	"net/url"

	"github.com/nikoksr/proji/pkg/storage/models"

	"github.com/nikoksr/proji/pkg/repo"
	"github.com/nikoksr/proji/pkg/util"
	"github.com/spf13/cobra"
)

var remoteRepos, directories, configs, excludes, packages, collections []string

const (
	flagCollection = "collection"
	flagConfig     = "config"
	flagDirectory  = "directory"
	flagExclude    = "exclude"
	flagPackage    = "package"
	flagRemoteRepo = "remote-repo"
)

var classImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more classes",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 && len(remoteRepos) < 1 && len(packages) < 1 && len(collections) < 1 {
			return fmt.Errorf("no flag given")
		}
		excludes = append(excludes, projiEnv.ExcludedPaths...)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Concat the two arrays so that '... import --config *.toml' is a valid command.
		// Without appending the args, proji would only use the first toml-file and not all of
		// them as intended with the '*'.
		configs = append(configs, args...)
		importTypes := map[string][]string{
			flagConfig:     configs,
			"dir":          directories,
			"repo":         remoteRepos,
			flagPackage:    packages,
			flagCollection: collections,
		}

		// Import configs
		for importType, paths := range importTypes {
			for _, path := range paths {
				result, err := importClass(path, importType, excludes)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Println(result)
				}
			}
		}
	},
}

func init() {
	classCmd.AddCommand(classImportCmd)

	classImportCmd.Flags().StringSliceVar(&remoteRepos, flagRemoteRepo, make([]string, 0), "create an importable config based on on the structure of a remote repository")
	_ = classImportCmd.MarkFlagDirname(flagRemoteRepo)

	classImportCmd.Flags().StringSliceVar(&directories, flagDirectory, make([]string, 0), "create an importable config based on the structure of a local directory")
	_ = classImportCmd.MarkFlagDirname(flagDirectory)

	classImportCmd.Flags().StringSliceVar(&configs, flagConfig, make([]string, 0), "import a class from a config file")
	_ = classImportCmd.MarkFlagFilename(flagConfig)

	classImportCmd.Flags().StringSliceVar(&packages, flagPackage, make([]string, 0), "import a package (EXPERIMENTAL)")
	_ = classImportCmd.MarkFlagFilename(flagPackage)

	classImportCmd.Flags().StringSliceVar(&collections, flagCollection, make([]string, 0), "import a collection of packages (EXPERIMENTAL)")
	_ = classImportCmd.MarkFlagFilename(flagCollection)

	classImportCmd.Flags().StringSliceVar(&excludes, flagExclude, make([]string, 0), "folder to exclude from local directory import")
	_ = classImportCmd.MarkFlagFilename(flagExclude)
}

func importClass(path, importType string, excludes []string) (string, error) {
	if util.IsInSlice(excludes, path) {
		return "", nil
	}

	class := models.NewClass("", "", false)
	var err error
	var confName, msg string
	var URL *url.URL
	var importer repo.Importer

	// In case of a repo, package or collection try to parse the path to a URL structure
	if importType == "repo" || importType == flagPackage || importType == flagCollection {
		URL, err = repo.ParseURL(path)
		if err != nil {
			return "", err
		}

		importer, err = models.GetRepoImporterFromURL(URL, projiEnv.Auth)
		if err != nil {
			return "", err
		}
	}

	switch importType {
	case flagConfig:
		err = class.ImportConfig(path)
		if err != nil {
			return "", err
		}
		err = projiEnv.StorageService.SaveClass(class)
		if err == nil {
			msg = fmt.Sprintf("> Successfully imported class '%s' from '%s'", class.Name, path)
		}
	case "dir":
		err = class.ImportFolderStructure(path, excludes)
		if err != nil {
			return "", err
		}
	case "repo":
		err = class.ImportRepoStructure(importer, nil)
		if err != nil {
			return "", err
		}
	case flagPackage:
		err = class.ImportPackage(URL, importer)
		if err != nil {
			return "", err
		}
		err = projiEnv.StorageService.SaveClass(class)
		if err == nil {
			msg = fmt.Sprintf("> Successfully imported class '%s' from '%s'", class.Name, path)
		}
	case flagCollection:
		classList, err := models.ImportClassesFromCollection(URL, importer)
		if err != nil {
			return "", err
		}
		for _, class := range classList {
			err = projiEnv.StorageService.SaveClass(class)
			if err == nil {
				msg += fmt.Sprintf("> Successfully imported class '%s' from '%s'\n", class.Name, path)
			} else {
				msg += fmt.Sprintf("> Importing class '%s' from '%s' failed: %v\n", class.Name, path, err)
			}
		}
	default:
		err = fmt.Errorf("path type %s is not supported", importType)
	}

	// Classes that are generated from directories or repos (structure, package and collection) should be exported to a config file first
	// so that the user can fine tune them
	if importType != flagConfig && importType != flagPackage && importType != flagCollection {
		confName, err = class.Export(".")
		if err == nil {
			msg = fmt.Sprintf("> '%s' was successfully exported to '%s'", path, confName)
		}
	}
	return msg, err
}
