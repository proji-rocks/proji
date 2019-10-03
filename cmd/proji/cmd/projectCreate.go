package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
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
			if err := createProject(name, label, cwd, projiEnv.ConfPath, projiEnv.Svc); err != nil {
				fmt.Printf("Creating project %s failed: %v\n", name, err)

				if err.Error() == "Project already exists" {
					if !helper.WantTo("Do you want to replace it?") {
						continue
					}
					if err := replaceProject(name, label, cwd, projiEnv.ConfPath, projiEnv.Svc); err != nil {
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

func createProject(name, label, cwd, configPath string, svc storage.Service) error {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}

	class, err := svc.LoadClass(classID)
	if err != nil {
		return err
	}

	// Load status active by default
	var status *item.Status
	status, err = svc.LoadStatus(1)
	if err != nil {
		return err
	}

	label = strings.ToLower(label)
	proj, err := item.NewProject(0, name, cwd+"/"+name, class, status)
	if err != nil {
		return err
	}

	// Save it first to see if it already exists in the database
	if err := svc.SaveProject(proj); err != nil {
		return err
	}
	// Create the project
	if err := proj.Create(cwd, configPath); err != nil {
		return err
	}
	return nil
}

func replaceProject(name, label, cwd, configPath string, svc storage.Service) error {
	id, err := svc.LoadProjectID(cwd + "/" + name)
	if err != nil {
		return err
	}

	// Replace it
	if err = svc.RemoveProject(id); err != nil {
		return err
	}
	return createProject(name, label, cwd, configPath, svc)
}
