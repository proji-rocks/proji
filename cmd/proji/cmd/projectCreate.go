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
	Use:   "create LABEL NAME [NAME...]",
	Short: "Create one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("At least one label and name have to be given")
		}
		label := args[0]
		projects := args[1:]

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		for _, name := range projects {
			if err := CreateProject(name, label, cwd); err != nil {
				fmt.Printf("Creating project %s failed: %v\n", name, err)
				if err.Error() == "project already exists" {
					if !wantToReplace() {
						continue
					}
					if err := replaceProject(name, label, cwd); err != nil {
						fmt.Printf("Replacing project %s failed: %v\n", name, err)
						continue
					}
					fmt.Printf("Project %s was successfully replaced.\n", name)
				}
				continue
			}
			fmt.Printf("Project %s was successfully created.\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// CreateProject will create a new project or return an error if the project already exists.
func CreateProject(name, label, cwd string) error {
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

	label = strings.ToLower(label)
	proj, err := storage.NewProject(name, label, cwd, s)
	if err != nil {
		return err
	}

	// Track it first to see if it already exists in the database
	if err := s.TrackProject(proj); err != nil {
		return err
	}
	// Create the project
	if err := proj.Create(); err != nil {
		return err
	}
	return nil
}

func wantToReplace() bool {
	// Ask to replace project
	var input string
	for {
		fmt.Print("Do you want to replace it? [y/n] ")
		n, err := fmt.Scan(&input)
		if n == 1 && err == nil {
			input = strings.Trim(strings.ToLower(input), " ")
			if input == "n" || input == "" {
				return false
			}
			if input == "y" {
				return true
			}
		}
	}
}

func replaceProject(name, label, cwd string) error {
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

	id, err := s.LoadProjectID(cwd + "/" + name)
	if err != nil {
		return err
	}

	// Replace it
	if err = s.UntrackProject(id); err != nil {
		return err
	}
	return CreateProject(name, label, cwd)
}
