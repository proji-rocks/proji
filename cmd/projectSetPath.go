package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectSetPathCmd = &cobra.Command{
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

		err = projiEnv.StorageService.UpdateProjectLocation(oldPath, newPath)
		if err != nil {
			fmt.Printf("> Setting path '%s' for project %s failed: %v\n", newPath, oldPath, err)
			return err
		}
		fmt.Printf("> Path '%s' was successfully set for project %s\n", newPath, oldPath)
		return nil
	},
}

func init() {
	projectSetCmd.AddCommand(projectSetPathCmd)
}
