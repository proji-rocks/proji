package cmd

import (
	"github.com/nikoksr/proji/pkg/util"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleanProjects()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func cleanProjects() error {
	projects, err := projiEnv.Svc.LoadAllProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		// Check path
		if util.DoesPathExist(project.Path) {
			continue
		}
		// Remove the project
		err := projiEnv.Svc.RemoveProject(project.Path)
		if err != nil {
			return err
		}
	}
	return nil
}