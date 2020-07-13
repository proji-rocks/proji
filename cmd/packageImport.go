package cmd

import (
	"fmt"
	"regexp"

	"github.com/nikoksr/proji/internal/util"

	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/spf13/cobra"
)

const (
	flagFilter             = "filter"
	flagConfig             = "config"
	flagDirectoryStructure = "dir-structure"
	flagRepoStructure      = "remote-structure"
	flagCollection         = "collection"
	flagPackage            = "package"
)

type packageImportCommand struct {
	cmd *cobra.Command
}

func newPackageImportCommand() *packageImportCommand {
	var remoteRepos, directories, configs, filters, packages, collections []string

	var cmd = &cobra.Command{
		Use:   "import LOCATION [LOCATION...]",
		Short: "Import one or more packages",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				if len(args) < 1 {
					return fmt.Errorf("no config path or flag given")
				}
				message.Warningf("no flag given, trying regular package import by default")
				packages = args
			} else {
				// Concat the two arrays so that '... import --config *.toml' is a valid command.
				// Without appending the args, proji would only use the first toml-file and not all of
				// them as intended with the '*'.
				configs = append(configs, args...)
			}
			filters = append(session.config.ExcludedPaths, filters...)
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

			// Cast filters
			regexFilters, err := util.StringsToRegex(filters)
			if err != nil {
				message.Errorf("failed to cast filters", err)
				return
			}

			// Import configs
			for importType, paths := range importTypes {
				for _, path := range paths {
					err := importPackage(path, importType, regexFilters)
					if err != nil {
						message.Warningf("failed to import package, %s", err.Error())
					}
				}
			}
		},
	}
	cmd.Flags().StringSliceVar(&packages, flagPackage, make([]string, 0), "import a package (default) (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&collections, flagCollection, "l", make([]string, 0), "import a collection of packages (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&configs, flagConfig, "c", make([]string, 0), "import a package from a config file")
	cmd.Flags().StringSliceVarP(&remoteRepos, flagRepoStructure, "r", make([]string, 0), "create an importable config based on on the structure of a remote repository")
	cmd.Flags().StringSliceVarP(&directories, flagDirectoryStructure, "d", make([]string, 0), "create an importable config based on the structure of a local directory")
	cmd.Flags().StringSliceVarP(&filters, flagFilter, "f", make([]string, 0), "filter imports with regex (only works with -l, -r, -d)")

	_ = cmd.MarkFlagDirname(flagDirectoryStructure)
	_ = cmd.MarkFlagFilename(flagConfig)

	return &packageImportCommand{cmd: cmd}
}

func importPackage(path, importType string, filters []*regexp.Regexp) error {
	var err error
	switch importType {
	case flagConfig:
		err = importPackageFromConfig(path)
	case flagDirectoryStructure:
		err = importPackageFromDirectoryStructure(path, filters)
	case flagRepoStructure:
		err = importPackageFromRepoStructure(path, filters)
	case flagCollection:
		err = importPackagesFromCollection(path, filters)
	case flagPackage:
		err = importPackageFromRemote(path)
	default:
		err = fmt.Errorf("import type not supported")
	}
	return err
}

func importPackageFromConfig(path string) error {
	// Import the package
	pkg, err := session.packageService.ImportPackageFromConfig(path)
	if err != nil {
		return err
	}

	// Save the package
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return err
	}
	message.Successf("successfully imported package %s", pkg.Name, path)
	return nil
}

func importPackageFromDirectoryStructure(path string, filters []*regexp.Regexp) error {
	// Import the package
	pkg, err := session.packageService.ImportPackageFromDirectoryStructure(path, filters)
	if err != nil {
		return err
	}

	// Save the package
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return err
	}
	message.Successf("successfully imported package %s", pkg.Name, path)
	return nil
}

func importPackageFromRepoStructure(url string, filters []*regexp.Regexp) error {
	// Parse url string to object
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		return err
	}

	// Import the package
	pkg, err := session.packageService.ImportPackageFromRepositoryStructure(parsedURL, filters)
	if err != nil {
		return err
	}

	// Save the package
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return err
	}
	message.Successf("successfully imported package %s", pkg.Name)
	return nil
}

func importPackagesFromCollection(url string, filters []*regexp.Regexp) error {
	// Parse url string to object
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		return err
	}

	// Import the packages
	pkgs, err := session.packageService.ImportPackagesFromCollection(parsedURL, filters)
	if err != nil {
		return err
	}

	// Save the packages
	for _, pkg := range pkgs {
		err = session.packageService.StorePackage(pkg)
		if err != nil {
			return err
		}
		message.Successf("successfully imported package %s", pkg.Name)
	}
	return nil
}

func importPackageFromRemote(url string) error {
	// Parse url string to object
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		return err
	}

	// Import the package
	pkg, err := session.packageService.ImportPackageFromRemote(parsedURL)
	if err != nil {
		return err
	}

	// Save the package
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return err
	}
	message.Successf("successfully imported package %s", pkg.Name)
	return nil
}
