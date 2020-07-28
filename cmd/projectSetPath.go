package cmd

import (
	"path/filepath"

	"github.com/nikoksr/proji/internal/message"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type projectSetPath struct {
	cmd *cobra.Command
}

func newProjectSetPathCommand() *projectSetPath {
	cmd := &cobra.Command{
		Use:                   "path OLD-PATH NEW-PATH",
		Short:                 "Set a new path",
		Aliases:               []string{"p"},
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			newPath, err := filepath.Abs(args[1])
			if err != nil {
				return err
			}

			err = session.projectService.UpdateProjectLocation(oldPath, newPath)
			if err != nil {
				return errors.Wrap(err, "failed setting project path")
			}
			message.Successf("successfully set path of project at %s to %s", oldPath, newPath)
			return nil
		},
	}
	return &projectSetPath{cmd: cmd}
}
