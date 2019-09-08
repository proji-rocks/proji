package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/project"
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
		for _, projectName := range projects {
			if err := project.CreateProject(label, projectName); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
