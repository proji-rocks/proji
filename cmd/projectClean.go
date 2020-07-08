//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"github.com/nikoksr/proji/messages"
	"github.com/nikoksr/proji/util"
	"github.com/pkg/errors"
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
	projects, err := session.StorageService.LoadProjects()
	if err != nil {
		return errors.Wrap(err, "failed to load all projects")
	}

	for _, project := range projects {
		// Check path
		if util.DoesPathExist(project.Path) {
			continue
		}
		// Remove the project
		err := session.StorageService.RemoveProject(project.Path)
		if err != nil {
			messages.Warning("failed to remove project with path %s, %s", project.Path, err.Error())
		}
	}
	return nil
}
