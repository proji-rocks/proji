package cmd

import (
	"bufio"
	"fmt"
	"github.com/nikoksr/proji/pkg/domain"
	packagestore "github.com/nikoksr/proji/pkg/package/store"
	"os"
	"regexp"

	"github.com/nikoksr/proji/internal/statuswriter"

	"github.com/nikoksr/proji/internal/message"
	"github.com/nikoksr/proji/pkg/remote"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	flagExclude            = "exclude"
	flagConfig             = "config"
	flagDirectoryStructure = "dir-structure"
	flagRepoStructure      = "remote-structure"
	flagCollection         = "collection"
	flagPackage            = "package"
	flagStdin              = "stdin"
)

type packageImportCommand struct {
	cmd *cobra.Command
}

func newPackageImportCommand() *packageImportCommand {
	var remoteRepos, directories, configs, packages, collections []string
	var fromStdin bool

	cmd := &cobra.Command{
		Use:     "import FROM [FROM...]",
		Short:   "Import one or more packages",
		Aliases: []string{"i"},
		Example: `  proji package import gh:nikoksr/proji-official-collection/configs/nikoksr/go.toml
  proji package import -r https://github.com/torvalds/linux
  proji package import -d .`,
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
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Compile exclude flag value to regex
			regexExclude, err := regexp.Compile(session.config.Import.Exclude)
			if err != nil {
				return errors.Wrap(err, "compile regex exclude")
			}

			// Associate different import types with their lists of config origins
			importTypesAndOrigins := map[string][]string{
				flagConfig:             configs,
				flagDirectoryStructure: directories,
				flagRepoStructure:      remoteRepos,
				flagPackage:            packages,
				flagCollection:         collections,
			}

			// Other imports have to be skipped in case of a stdin import because they might pollute stdin and
			// potentially make the package import fail.
			if fromStdin {
				importTypesAndOrigins = map[string][]string{flagStdin: {""}}
			}

			// Loop over all import types and try to import packages from their associated config origins
			sw := statuswriter.New()
			sw.Run()
			for importType, origins := range importTypesAndOrigins {
				for _, origin := range origins {
					go importPackage(sw.NewSink(), origin, importType, regexExclude)
				}
			}
			sw.Wait()

			return nil
		},
	}
	cmd.Flags().StringSliceVar(&packages, flagPackage, make([]string, 0), "import a package (default) (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&collections, flagCollection, "c", make([]string, 0), "import a collection of packages (EXPERIMENTAL)")
	cmd.Flags().StringSliceVarP(&configs, flagConfig, "f", make([]string, 0), "import a package from a config file")
	cmd.Flags().StringSliceVarP(&remoteRepos, flagRepoStructure, "r", make([]string, 0), "create an importable config based on on the structure of a remote repository")
	cmd.Flags().StringSliceVarP(&directories, flagDirectoryStructure, "d", make([]string, 0), "create an importable config based on the structure of a local directory")
	cmd.Flags().StringP(flagExclude, "e", "", "regex pattern to exclude paths from import (only works with -c, -r, -d)")
	cmd.Flags().BoolVarP(&fromStdin, flagStdin, "i", false, "import a configuration from standard input")

	_ = cmd.MarkFlagDirname(flagDirectoryStructure)
	_ = cmd.MarkFlagFilename(flagConfig)

	return &packageImportCommand{cmd: cmd}
}

func importPackage(status *statuswriter.Sink, path, importType string, exclude *regexp.Regexp) {
	defer status.Close()
	var pkg *domain.Package
	var err error

	status.Write(message.Sinfof("importing %s %s\n", importType, path))
	switch importType {
	case flagConfig:
		pkg, err = importPackageFromConfig(status, path)
	case flagDirectoryStructure:
		pkg, err = importPackageFromDirectoryStructure(status, path, exclude)
	case flagRepoStructure:
		pkg, err = importPackageFromRepoStructure(status, path, exclude)
	case flagPackage:
		pkg, err = importPackageFromRemote(status, path)
	case flagCollection:
		importPackagesFromCollection(status, path, exclude)
	case flagStdin:
		pkg, err = importPackageFromStdin(status)
	default:
		err = fmt.Errorf("import type %s not supported", importType)
	}
	if errors.Is(err, packagestore.ErrPackageExists) && importType != flagConfig {
		handleDuplicatePackage(status, pkg)
		return
	}
	if err != nil {
		status.Write(message.Serrorf(err, "failed to import package from %s %s", importType, path))
		return
	}
	if importType != flagCollection {
		// Collections handle messages on its own
		status.Write(message.Ssuccessf("successfully imported package %s [%s]", pkg.Name, pkg.Label))
	}
}

func importPackageFromConfig(status *statuswriter.Sink, path string) (*domain.Package, error) {
	// Import the package
	pkg, err := session.packageService.ImportPackageFromConfig(path)
	if err != nil {
		return nil, err
	}

	// Save the package
	status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
	err = session.packageService.StorePackage(pkg)
	return pkg, err
}

func importPackageFromDirectoryStructure(status *statuswriter.Sink, path string, exclude *regexp.Regexp) (*domain.Package, error) {
	// Import the package
	pkg, err := session.packageService.ImportPackageFromDirectoryStructure(path, exclude)
	if err != nil {
		return nil, err
	}

	// Save the package
	status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return nil, err
	}
	return pkg, err
}

func importPackageFromRepoStructure(status *statuswriter.Sink, url string, exclude *regexp.Regexp) (*domain.Package, error) {
	// Parse url string to object
	status.Write(message.Sinfof("parsing url"))
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// Import the package
	status.Write(message.Sinfof("creating package from repository structure of %s", parsedURL.String()))
	pkg, err := session.packageService.ImportPackageFromRepositoryStructure(parsedURL, exclude)
	if err != nil {
		return nil, err
	}

	// Save the package
	status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
	err = session.packageService.StorePackage(pkg)
	return pkg, err
}

func importPackageFromRemote(status *statuswriter.Sink, url string) (*domain.Package, error) {
	// Parse url string to object
	status.Write(message.Sinfof("parsing url"))
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// Import the package
	status.Write(message.Sinfof("importing package from %s", parsedURL.String()))
	pkg, err := session.packageService.ImportPackageFromRemote(parsedURL)
	if err != nil {
		return nil, err
	}

	// Save the package
	status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
	err = session.packageService.StorePackage(pkg)
	if err != nil {
		return nil, err
	}
	message.Successf("successfully imported package %s", pkg.Name)
	return pkg, err
}

func importPackagesFromCollection(status *statuswriter.Sink, url string, exclude *regexp.Regexp) {
	// Parse url string to object
	status.Write(message.Sinfof("parsing url"))
	parsedURL, err := remote.ParseURL(url)
	if err != nil {
		status.Write(message.Serrorf(err, "failed to parse collection url %s", url))
		return
	}

	// Import the packages
	status.Write(message.Sinfof("importing packages from collection %s", parsedURL.String()))
	pkgs, err := session.packageService.ImportPackagesFromCollection(parsedURL, exclude)
	if err != nil {
		status.Write(message.Serrorf(err, "failed to import packages from collection %s", parsedURL.String()))
		return
	}

	// Save the packages
	var successfulImports int
	for _, pkg := range pkgs {
		status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
		err = session.packageService.StorePackage(pkg)
		if errors.Is(err, packagestore.ErrPackageExists) {
			handleDuplicatePackage(status, pkg)
			continue
		}
		if err != nil {
			status.Write(message.Serrorf(err, "failed to store package %s [%s]", pkg.Name, pkg.Label))
		} else {
			status.Write(message.Ssuccessf("successfully imported package %s [%s]", pkg.Name, pkg.Label))
			successfulImports++
		}
	}
	status.Write(message.Ssuccessf("successfully imported %d of %d package from collection %s", successfulImports, len(pkgs), parsedURL.String()))
}

func importPackageFromStdin(status *statuswriter.Sink) (*domain.Package, error) {
	// Scan package config from stdin into string
	status.Write(message.Sinfof("importing package from stdin"))
	scanner := bufio.NewScanner(os.Stdin)
	input := ""
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}

	// Import actual package from the string
	pkg, err := session.packageService.ImportPackageFromString(input)
	if err != nil {
		return nil, err
	}

	// Save the package
	status.Write(message.Sinfof("storing package %s [%s]", pkg.Name, pkg.Label))
	err = session.packageService.StorePackage(pkg)

	return pkg, err
}

func handleDuplicatePackage(status *statuswriter.Sink, pkg *domain.Package) {
	// Announce config export
	status.Write(message.Swarningf(
		"%v (label %s): exporting package config",
		packagestore.ErrPackageExists,
		pkg.Label),
	)

	// Try to export package config for editing
	exportedTo, err := session.packageService.ExportPackageToConfig(*pkg, ".")
	if err != nil {
		status.Write(message.Serrorf(
			err,
			"%v (label %s): failed to export package config",
			packagestore.ErrPackageExists,
			pkg.Label),
		)
		return
	}
	status.Write(message.Swarningf(
		"%v (label %s): exported package config to %s",
		packagestore.ErrPackageExists,
		pkg.Label,
		exportedTo),
	)
}
