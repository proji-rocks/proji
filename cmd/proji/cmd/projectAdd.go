package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add LABEL PATH STATUS",
	Short: "Add an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("Missing label, path or status")
		}

		path, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		if !helper.DoesPathExist(path) {
			return fmt.Errorf("path '%s' does not exist", path)
		}

		label := strings.ToLower(args[0])
		status := strings.ToLower(args[2])

		if err := addProject(label, path, status); err != nil {
			return err
		}
		fmt.Printf("Project '%s' was successfully added.\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addProject(label, path, statusTitle string) error {
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

	name := filepath.Base(path)
	if err != nil {
		return err
	}

	classID, err := s.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}

	statusID, err := s.LoadStatusID(statusTitle)
	if err != nil {
		return err
	}

	proj, err := storage.NewProject(0, name, path, classID, statusID, s)
	if err != nil {
		return err
	}
	if err := s.SaveProject(proj); err != nil {
		return err
	}
	return nil
}
