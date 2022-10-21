package pkg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/nikoksr/proji/internal/cli"
	"github.com/nikoksr/proji/internal/text"
	"github.com/nikoksr/proji/pkg/logging"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ls",
		Short:                 "List installed packages",
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return listPackages(cmd.Context())
		},
	}

	return cmd
}

func listPackages(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	// Get package manager from session
	logger.Debug("getting package manager from cli session")
	pama := cli.SessionFromContext(ctx).PackageManager
	if pama == nil {
		return errors.New("no package manager available")
	}

	// Call the packageList.
	logger.Debug("fetching package list")
	packageList, err := pama.Fetch(ctx)
	if err != nil {
		return errors.Wrap(err, "fetch packages")
	}

	// Exit when no packages are installed.
	if len(packageList) == 0 {
		return nil
	}

	// List all packageList in a pretty table.
	logger.Debug("listing packages")

	table := text.NewTablePrinter()
	table.AddHeaderColumns("#", "Label", "Name", "Upstream", "Created at")

	for idx, _package := range packageList {
		table.AddRow(idx+1, _package.Label, _package.Name, _package.UpstreamURL, _package.CreatedAt)
	}

	err = table.Render()
	if err != nil {
		return errors.Wrap(err, "render table")
	}

	return nil
}
