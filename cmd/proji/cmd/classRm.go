package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/proji/storage"
	"github.com/spf13/cobra"
)

var removeAll bool

var classRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if removeAll {
			err := removeAllClasses(projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing of all classes failed: %v\n", err)
				return err
			}
			fmt.Println("> All classes were successfully removed")
			return nil
		}

		if len(args) < 1 {
			return fmt.Errorf("missing class label")
		}

		for _, name := range args {
			err := removeClass(name, projiEnv.Svc)
			if err != nil {
				fmt.Printf("> Removing '%s' failed: %v\n", name, err)
				continue
			}
			fmt.Printf("> '%s' was successfully removed\n", name)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
	classRmCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all classes")
}

func removeClass(label string, svc storage.Service) error {
	classID, err := svc.LoadClassIDByLabel(label)
	if err != nil {
		return err
	}
	class, err := svc.LoadClass(classID)
	if err != nil {
		return err
	}

	if class.IsDefault {
		return fmt.Errorf("default classes can not be removed")
	}

	return svc.RemoveClass(classID)
}

func removeAllClasses(svc storage.Service) error {
	classes, err := svc.LoadAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		if class.IsDefault {
			continue
		}
		err = svc.RemoveClass(class.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
