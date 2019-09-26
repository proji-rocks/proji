package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add LABEL PATH",
	Short: "Add an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("Missing label or path")
		}

		path, err := filepath.Abs(args[1])
		if !helper.DoesPathExist(path) {
			return fmt.Errorf("path '%s' does not exist", path)
		}

		name := filepath.Base(path)
		if err != nil {
			return err
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		label := strings.ToLower(args[0])

		if err := AddProject(name, label, cwd); err != nil {
			return err
		}
		fmt.Printf("Project '%s' was successfully added.\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// AddProject will create a new project or return an error if the project already exists.
func AddProject(name, label, cwd string) error {
	// Setup storage
	sqlitePath, err := helper.GetSqlitePath()
	if err != nil {
		return err
	}
	s, err := sqlite.New(sqlitePath)
	if err != nil {
		return err
	}
	defer s.Close()

	proj, err := storage.NewProject(name, label, cwd, s)
	if err != nil {
		return err
	}
	if err := s.TrackProject(proj); err != nil {
		return err
	}
	return nil
}
