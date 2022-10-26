package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nikoksr/simplog"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/pkg/api/v1/domain"
)

func newShowCommand() *cobra.Command {
	var showAll bool

	cmd := &cobra.Command{
		Use:                   "show [OPTIONS] LABEL [LABEL...]",
		Short:                 "Show details about installed packages",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			if !showAll && len(args) < 1 {
				return fmt.Errorf("missing package label")
			}

			var labels []string
			if !showAll {
				labels = args
			}

			return showPackages(cmd.Context(), labels...)
		},
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all packages")

	return cmd
}

func prettyPrintPackage(_ context.Context, _package *domain.Package) error {
	pkgJSON, err := json.MarshalIndent(_package, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal package")
	}

	fmt.Println(string(pkgJSON))

	return nil
}

func showPackages(ctx context.Context, labels ...string) error {
	logger := simplog.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Showing packages
	for _, label := range labels {
		logger.Debugf("showing package %q", label)
		pkg, err := pama.GetByLabel(ctx, label)
		if err != nil {
			return errors.Wrapf(err, "get package %q", label)
		}

		logger.Debugf("pretty printing package %q", label)
		if err = prettyPrintPackage(ctx, &pkg); err != nil {
			return errors.Wrapf(err, "pretty print package %q", label)
		}
	}

	return nil
}
