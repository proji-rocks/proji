package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/nikoksr/proji/messages"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type projectSetPath struct {
	cmd *cobra.Command
}

func newProjectSetPathCommand() *projectSetPath {
	var cmd = &cobra.Command{
		Use:   "path OLD-PATH NEW-PATH",
		Short: "Set a new path",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing old or new path")
			}

			oldPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			newPath, err := filepath.Abs(args[1])
			if err != nil {
				return err
			}

			err = activeSession.storageService.UpdateProjectLocation(oldPath, newPath)
			if err != nil {
				return errors.Wrap(err, "failed setting project path")
			}
			messages.Success("successfully set path of project at %s to %s", oldPath, newPath)
			return nil
		},
	}
	return &projectSetPath{cmd: cmd}
}
