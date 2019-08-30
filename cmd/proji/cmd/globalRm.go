package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/internal/app/proji/global"
	"github.com/spf13/cobra"
)

var globalRmType string

var globalRmCmd = &cobra.Command{
	Use:   "rm GLOBAL-ID [GLOBAL-ID...]",
	Short: "remove existing globals by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing global id")
		}
		if err := isTypeValid(globalRmType); err != nil {
			return err
		}
		return global.RemoveGlobal(globalRmType, args)
	},
}

func init() {
	globalCmd.AddCommand(globalRmCmd)

	// Flag to define type of global
	globalRmCmd.PersistentFlags().StringVarP(&globalRmType, "type", "t", "", "type of global - folder, file or script")
	globalRmCmd.MarkPersistentFlagRequired("type")
}
