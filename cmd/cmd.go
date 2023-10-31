package cmd

import (
	"io"
	"log"

	"github.com/spf13/cobra"
)

var (
	rootOpts struct {
		verbose bool
	}

	rootCmd = &cobra.Command{
		Use:   "sqlgen",
		Short: "`sqlgen` is a Go SQL code generator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootOpts.verbose {
				log.SetFlags(0)
			} else {
				log.SetOutput(io.Discard)
			}
			return cmd.Help()
		},
	}
)

func Execute() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.verbose, "verbose", "v", false, "Shows the log details.")
	cobra.CheckErr(rootCmd.Execute())
}
