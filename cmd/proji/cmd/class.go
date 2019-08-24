package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var classCmd = &cobra.Command{
	Use:   "class NAME",
	Short: "Add a new class to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Missing class name")
		}

		fmt.Println("Added new class " + args[0])
		return nil
	},
}

func init() {
	addCmd.AddCommand(classCmd)
}
