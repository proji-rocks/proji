package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var projectSetPathCmd = &cobra.Command{
	Use:   "path PATH PROJECT-ID",
	Short: "Set a new path",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Missing path or project-ID")
		}

		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		projectID, err := helper.StrToUInt(args[1])
		if err != nil {
			return err
		}

		if err := setPath(projectID, path); err != nil {
			fmt.Printf("Setting path '%s' for project %d failed: %v\n", path, projectID, err)
			return err
		}
		fmt.Printf("Path '%s' was successfully set for project %d.\n", path, projectID)
		return nil
	},
}

func init() {
	projectSetCmd.AddCommand(projectSetPathCmd)
}

func setPath(projectID uint, path string) error {
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.UpdateProjectLocation(projectID, path)
}
