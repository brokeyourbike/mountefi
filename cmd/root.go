package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mountefi",
	Short: "mountefi - an EFI Mounting Utility",
	Long: `mountefi - an EFI Mounting Utility with no dependencies.
Documentation is available at https://github.com/brokeyourbike/mountefi`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
