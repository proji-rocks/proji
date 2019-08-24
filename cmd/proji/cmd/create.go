package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/project"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create EXTENSION PROJECT [PROJECTS]",
	Short: "Create new projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("insufficient number of cli arguments")
		}

		project.CreateProject(args[0], args[1:])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
