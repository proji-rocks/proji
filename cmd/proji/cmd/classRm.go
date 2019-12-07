package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/spf13/cobra"
)

var removeAll bool

var classRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect classes that will be removed
		var classes []*item.Class

		if removeAll {
			var err error
			classes, err = projiEnv.Svc.LoadAllClasses()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing class label")
			}

			for _, label := range args {
				classID, err := projiEnv.Svc.LoadClassIDByLabel(label)
				if err != nil {
					return err
				}
				class, err := projiEnv.Svc.LoadClass(classID)
				if err != nil {
					return err
				}
				classes = append(classes, class)
			}
		}

		// Remove the classes
		for _, class := range classes {
			if class.IsDefault {
				continue
			}
			err := projiEnv.Svc.RemoveClass(class.ID)
			if err != nil {
				fmt.Printf("> Removing '%s' failed: %v\n", class.Label, err)
				return err
			}
			fmt.Printf("> '%s' was successfully removed\n", class.Label)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
	classRmCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all classes")
}
