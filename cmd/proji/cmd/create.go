package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/proji/storage"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create LABEL PROJECTNAME [PROJECTNAME...]",
	Short: "create new projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("atleast one project name has to be specified")
		}
		label := args[0]
		projects := args[1:]
		for _, name := range projects {
			if err := CreateProject(name, label); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// CreateProject will create a new project or return an error if the project already exists.
func CreateProject(name, label string) error {
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

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Header
	fmt.Println(helper.ProjectHeader(name))

	label = strings.ToLower(label)
	proj, err := storage.NewProject(name, label, cwd, s)
	if err != nil {
		return err
	}

	// Create
	if err := proj.Create(); err != nil {
		return fmt.Errorf("could not create project %s: %v", proj.Name, err)
	}
	// Track
	if err := s.TrackProject(proj); err != nil {
		return fmt.Errorf("could not track project %s: %v", proj.Name, err)
	}
	return nil
}
