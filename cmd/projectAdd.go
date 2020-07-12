package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/pkg/domain"

	"github.com/nikoksr/proji/internal/message"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/internal/util"
	"github.com/spf13/cobra"
)

type projectAddCommand struct {
	cmd *cobra.Command
}

func newProjectAddCommand() *projectAddCommand {
	var cmd = &cobra.Command{
		Use:                   "add LABEL PATH",
		Short:                 "Add an existing project",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[1])
			if err != nil {
				return err
			}
			if !util.DoesPathExist(path) {
				return fmt.Errorf("path %s does not exist", path)
			}

			label := strings.ToLower(args[0])

			err = addProject(label, path)
			if err != nil {
				return errors.Wrap(err, "failed to add project")
			}
			message.Successf("successfully added project at path %s", path)
			return nil
		},
	}
	return &projectAddCommand{cmd: cmd}
}

func addProject(label, path string) error {
	name := filepath.Base(path)
	pkg, err := session.packageService.LoadPackage(label)
	if err != nil {
		return errors.Wrap(err, "failed to load package")
	}

	project := domain.NewProject(name, path, pkg)
	err = session.projectService.StoreProject(project)
	if err != nil {
		return errors.Wrap(err, "failed to save package")
	}
	return nil
}
