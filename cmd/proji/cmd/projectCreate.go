package cmd

import (
	"fmt"
	"os"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/nikoksr/proji/pkg/proji/storage/item"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create LABEL NAME [NAME...]",
	Short: "Create one or more projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("at least one label and name have to be given")
		}
		label := args[0]
		projects := args[1:]

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Load class and status once for all projects
		classID, err := projiEnv.Svc.LoadClassIDByLabel(label)
		if err != nil {
			return err
		}

		class, err := projiEnv.Svc.LoadClass(classID)
		if err != nil {
			return err
		}

		// Load status active by default
		status, err := projiEnv.Svc.LoadStatus(1)
		if err != nil {
			return err
		}

		for _, name := range projects {
			fmt.Printf("\n> Creating project %s\n", name)

			err := createProject(name, cwd, projiEnv.UserConfigPath, class, status)
			if err != nil {
				fmt.Printf(" -> Failed: %v\n", err)

				if err.Error() == "Project already exists" {
					if !helper.WantTo("> Do you want to replace it?") {
						continue
					}
					err := replaceProject(name, cwd, projiEnv.UserConfigPath, class, status)
					if err != nil {
						fmt.Printf("> Replacing project %s failed: %v\n", name, err)
						continue
					}
					fmt.Printf("> Project %s was successfully replaced\n", name)
				}
				continue
			}
			fmt.Printf("> Project %s was successfully created\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createProject(name, cwd, configPath string, class *item.Class, status *item.Status) error {
	proj := item.NewProject(0, name, cwd+"/"+name, class, status)

	// Save it first to see if it already exists in the database
	err := projiEnv.Svc.SaveProject(proj)
	if err != nil {
		return err
	}
	// Create the project
	err = proj.Create(cwd, configPath)
	if err != nil {
		return err
	}
	return nil
}

func replaceProject(name, cwd, configPath string, class *item.Class, status *item.Status) error {
	id, err := projiEnv.Svc.LoadProjectID(cwd + "/" + name)
	if err != nil {
		return err
	}

	// Replace it
	err = projiEnv.Svc.RemoveProject(id)
	if err != nil {
		return err
	}
	return createProject(name, cwd, configPath, class, status)
}
