package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unmountCmd)
}

var unmountCmd = &cobra.Command{
	Use:   "unmount",
	Short: "Unmount",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
