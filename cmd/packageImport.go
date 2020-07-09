//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"net/url"

	"github.com/nikoksr/proji/messages"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/storage/models"

	"github.com/nikoksr/proji/repo"
	"github.com/nikoksr/proji/util"
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

var packageImportCmd = &cobra.Command{
	Use:   "import FILE [FILE...]",
	Short: "Import one or more packages",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(configs) < 1 && len(directories) < 1 && len(remoteRepos) < 1 && len(packages) < 1 && len(collections) < 1 {
			return fmt.Errorf("no flag given")
		}
		excludes = append(excludes, session.Config.ExcludedPaths...)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Concat the two arrays so that '... import --config *.toml' is a valid command.
		// Without appending the args, proji would only use the first toml-file and not all of
		// them as intended with the '*'.
		configs = append(configs, args...)
		importTypes := map[string][]string{
			flagConfig:     configs,
			flagDirectory:  directories,
			flagRemoteRepo: remoteRepos,
			flagPackage:    packages,
			flagCollection: collections,
		}

		// Import configs
		for importType, paths := range importTypes {
			for _, path := range paths {
				result, err := importPackage(path, importType, excludes)
				if err != nil {
					messages.Warning("failed to import package, %s", err.Error())
				} else {
					messages.Success(result)
				}
			}
		}
	},
}

func init() {
	packageCmd.AddCommand(packageImportCmd)

	packageImportCmd.Flags().StringSliceVarP(&directories, flagDirectory, "d", make([]string, 0), "create an importable config based on the structure of a local directory")
	_ = packageImportCmd.MarkFlagDirname(flagDirectory)

	packageImportCmd.Flags().StringSliceVarP(&configs, flagConfig, "f", make([]string, 0), "import a package from a config file")
	_ = packageImportCmd.MarkFlagFilename(flagConfig)

	packageImportCmd.Flags().StringSliceVarP(&excludes, flagExclude, "e", make([]string, 0), "folder to exclude from local directory import")
	_ = packageImportCmd.MarkFlagFilename(flagExclude)

	packageImportCmd.Flags().StringSliceVarP(&remoteRepos, flagRemoteRepo, "r", make([]string, 0), "create an importable config based on on the structure of a remote repository")

	packageImportCmd.Flags().StringSliceVarP(&packages, flagPackage, "p", make([]string, 0), "import a package (EXPERIMENTAL)")

	packageImportCmd.Flags().StringSliceVarP(&collections, flagCollection, "c", make([]string, 0), "import a collection of packages (EXPERIMENTAL)")

}

func importPackage(path, importType string, excludes []string) (string, error) {
	if util.IsInSlice(excludes, path) {
		return "", nil
	}

	pkg := models.NewPackage("", "", false)
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

		importer, err = models.GetRepoImporterFromURL(URL, session.Config.Auth)
		if err != nil {
			return "", err
		}
	}

	switch importType {
	case flagConfig:
		err = pkg.ImportFromConfig(path)
		if err != nil {
			return "", err
		}
		err = session.StorageService.SavePackage(pkg)
		if err == nil {
			msg = fmt.Sprintf("successfully imported package %s from %s", pkg.Name, path)
		}
	case "dir":
		err = pkg.ImportFromFolderStructure(path, excludes)
		if err != nil {
			return "", err
		}
	case "repo":
		err = pkg.ImportFromRepoStructure(importer, nil)
		if err != nil {
			return "", errors.Wrap(err, "failed to import repository structure")
		}
	case flagPackage:
		err = pkg.ImportFromRepo(URL, importer)
		if err != nil {
			return "", errors.Wrap(err, "failed to import package from repository")
		}
		err = session.StorageService.SavePackage(pkg)
		if err != nil {
			return "", errors.Wrap(err, "failed to save package")
		}
		msg = fmt.Sprintf("successfully imported package %s from %s", pkg.Name, path)
	case flagCollection:
		packageList, err := models.ImportCollectionFromRepo(URL, importer)
		if err != nil {
			return "", errors.Wrap(err, "failed to import collection")
		}
		for _, pkg := range packageList {
			err = session.StorageService.SavePackage(pkg)
			if err != nil {
				msg += fmt.Sprintf("failed to import package %s from %s, %s", pkg.Name, path, err.Error())
			} else {
				msg += fmt.Sprintf("successfully imported package %s from %s", pkg.Name, path)
			}
		}
	default:
		err = fmt.Errorf("path type %s is not supported", importType)
	}

	// Packages that are generated from directories or repos (structure, package and collection) should be exported to a config file first
	// so that the user can fine tune them
	if importType != flagConfig && importType != flagPackage && importType != flagCollection {
		confName, err = pkg.ExportConfig(".")
		if err == nil {
			msg = fmt.Sprintf("successfully exported %s to %s", path, confName)
		}
	}
	return msg, err
}
