package cmd

import (
	"fmt"
	"net/url"

	"github.com/nikoksr/proji/messages"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"

	"github.com/nikoksr/proji/repo"
	"github.com/spf13/cobra"
)

const (
	flagExclude            = "exclude"
	flagConfig             = "config"
	flagDirectoryStructure = "dir-structure"
	flagRepoStructure      = "repo-structure"
	flagCollection         = "collection"
	flagPackage            = "package"
)

type packageImportCommand struct {
	cmd *cobra.Command
}

func newPackageImportCommand() *packageImportCommand {
	var remoteRepos, directories, configs, excludes, packages, collections []string

	var cmd = &cobra.Command{
		Use:   "import FILE [FILE...]",
		Short: "Import one or more packages",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				if len(args) < 1 {
					return fmt.Errorf("no config path or flag given")
				}
				messages.Warningf("no flag given, trying regular package import by default")
				packages = args
			} else {
				// Concat the two arrays so that '... import --config *.toml' is a valid command.
				// Without appending the args, proji would only use the first toml-file and not all of
				// them as intended with the '*'.
				configs = append(configs, args...)
			}
			excludes = append(activeSession.config.ExcludedPaths, excludes...)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			importTypes := map[string][]string{
				flagConfig:             configs,
				flagDirectoryStructure: directories,
				flagRepoStructure:      remoteRepos,
				flagPackage:            packages,
				flagCollection:         collections,
			}

			// Import configs
			for importType, paths := range importTypes {
				for _, path := range paths {
					err := importPackage(path, importType, excludes)
					if err != nil {
						messages.Warningf("failed to import package, %s", err.Error())
					}
				}
			}
		},
	}
	cmd.Flags().StringSliceVar(&packages, flagPackage, make([]string, 0), "import a package (default) (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&collections, flagCollection, "c", make([]string, 0), "import a collection of packages (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&configs, flagConfig, "f", make([]string, 0), "import a package from a config file")
	cmd.Flags().StringSliceVarP(&remoteRepos, flagRepoStructure, "r", make([]string, 0), "create an importable config based on on the structure of a remote repository")
	cmd.Flags().StringSliceVarP(&directories, flagDirectoryStructure, "d", make([]string, 0), "create an importable config based on the structure of a local directory")
	cmd.Flags().StringSliceVarP(&excludes, flagExclude, "e", make([]string, 0), "folder to exclude from local directory import")

	_ = cmd.MarkFlagDirname(flagDirectoryStructure)
	_ = cmd.MarkFlagFilename(flagConfig)
	_ = cmd.MarkFlagFilename(flagExclude)

	return &packageImportCommand{cmd: cmd}
}

func importPackage(path, importType string, excludes []string) error {
	var err error
	switch importType {
	case flagConfig:
		err = importPackageFromConfig(path)
	case flagDirectoryStructure:
		err = importPackageFromDirectoryStructure(path, excludes)
	case flagRepoStructure:
		err = importPackageFromRepoStructure(path)
	case flagCollection:
		err = importPackagesFromCollection(path)
	case flagPackage:
		err = importPackageFromRepo(path)
	}
	return err
}

func exportPackageConfig(pkg *models.Package) error {
	// Export package config to current working directory
	confName, err := pkg.ExportConfig(".")
	if err != nil {
		return err
	}
	messages.Successf("successfully imported package %s to %s", pkg.Name, confName)
	return nil
}

func importPackageFromConfig(path string) error {
	// Import the package
	pkg := models.NewPackage("", "", false)
	err := pkg.ImportFromConfig(path)
	if err != nil {
		return err
	}

	// Save the package
	err = activeSession.storageService.SavePackage(pkg)
	if err != nil {
		return err
	}
	messages.Successf("successfully imported package %s from %s", pkg.Name, path)
	return nil
}

func importPackageFromDirectoryStructure(path string, excludes []string) error {
	// Import the package
	pkg := models.NewPackage("", "", false)
	err := pkg.ImportFromFolderStructure(path, excludes)
	if err != nil {
		return err
	}
	// Export the config for user editing
	return exportPackageConfig(pkg)
}

func importPackageFromRepoStructure(url string) error {
	// Get repo importer
	_, importer, err := getURLAndRepoImporter(url)
	if err != nil {
		return err
	}

	// Import the package
	pkg := models.NewPackage("", "", false)
	err = pkg.ImportFromRepoStructure(importer, nil)
	if err != nil {
		return errors.Wrap(err, "import repository structure")
	}

	// Export the config for user editing
	return exportPackageConfig(pkg)
}

func importPackagesFromCollection(url string) error {
	// Get parsed url and repo importer
	parsedURL, importer, err := getURLAndRepoImporter(url)
	if err != nil {
		return err
	}

	// Import the packages
	packageList, err := models.ImportCollectionFromRepo(parsedURL, importer)
	if err != nil {
		return errors.Wrap(err, "failed to import collection")
	}

	// Save the packages to storage
	for _, pkg := range packageList {
		err = activeSession.storageService.SavePackage(pkg)
		if err != nil {
			messages.Warningf("failed to import package %s, %s", pkg.Name, err.Error())
		} else {
			messages.Successf("successfully imported package %s from %s", pkg.Name, importer.URL().String())
		}
	}
	return nil
}

func importPackageFromRepo(url string) error {
	// Get parsed url and repo importer
	parsedURL, importer, err := getURLAndRepoImporter(url)
	if err != nil {
		return err
	}

	// Import the package
	pkg := models.NewPackage("", "", false)
	err = pkg.ImportFromRepo(parsedURL, importer)
	if err != nil {
		return errors.Wrap(err, "failed to import package from repository")
	}

	// Save the package
	err = activeSession.storageService.SavePackage(pkg)
	if err != nil {
		return errors.Wrap(err, "failed to save package")
	}
	messages.Successf("successfully imported package %s from %s", pkg.Name, parsedURL.String())
	return nil
}

// getParsedURL tries to parse an url string to an url object. The parsing will validate the given url
// that way. If the url is valid, it returns a url object.
func getParsedURL(url string) (*url.URL, error) {
	parsedURL, err := repo.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	return parsedURL, nil
}

// getURLAndRepoImporter tries to parse a given url string to a url object. If successful (if url valid),
// it will try to get an importer interface for the given url. If this is successful too, the function
// will return the parsed url and the importer interface.
func getURLAndRepoImporter(url string) (*url.URL, repo.Importer, error) {
	parsedURL, err := getParsedURL(url)
	if err != nil {
		return nil, nil, err
	}
	importer, err := models.GetRepoImporterFromURL(parsedURL, activeSession.config.Auth)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get repo importer")
	}
	return parsedURL, importer, nil
}
