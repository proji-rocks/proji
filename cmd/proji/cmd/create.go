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
			return fmt.Errorf("insufficient number of cli arguments")
		}
		return project.CreateProject(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
