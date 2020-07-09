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

const (
	flagCollection = "collection"
	flagConfig     = "config"
	flagDirectory  = "directory"
	flagExclude    = "exclude"
	flagPackage    = "package"
	flagRemoteRepo = "remote-repo"
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
				messages.Warningf("no flag given, using --config by default")
				configs = args
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
						messages.Warningf("failed to import package, %s", err.Error())
					} else {
						messages.Successf(result)
					}
				}
			}
		},
	}
	cmd.Flags().StringSliceVarP(&directories, flagDirectory, "d", make([]string, 0), "create an importable config based on the structure of a local directory")
	cmd.Flags().StringSliceVarP(&configs, flagConfig, "f", make([]string, 0), "import a package from a config file")
	cmd.Flags().StringSliceVarP(&excludes, flagExclude, "e", make([]string, 0), "folder to exclude from local directory import")
	cmd.Flags().StringSliceVarP(&remoteRepos, flagRemoteRepo, "r", make([]string, 0), "create an importable config based on on the structure of a remote repository")
	cmd.Flags().StringSliceVarP(&packages, flagPackage, "p", make([]string, 0), "import a package (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&collections, flagCollection, "c", make([]string, 0), "import a collection of packages (EXPERIMENTAL)")

	_ = cmd.MarkFlagDirname(flagDirectory)
	_ = cmd.MarkFlagFilename(flagConfig)
	_ = cmd.MarkFlagFilename(flagExclude)

	return &packageImportCommand{cmd: cmd}
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

		importer, err = models.GetRepoImporterFromURL(URL, activeSession.config.Auth)
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
		err = activeSession.storageService.SavePackage(pkg)
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
		err = activeSession.storageService.SavePackage(pkg)
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
			err = activeSession.storageService.SavePackage(pkg)
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
