package cmd

import (
	"io"
	"log"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/spf13/cobra"
)

var (
	rootOpts struct {
		verbose bool
	}

	rootCmd = &cobra.Command{
		Use:   "sqlgen",
		Short: "🚀 Transform your struct to SQL Go code!!!",
		Long:  ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if rootOpts.verbose {
				log.SetFlags(0)
			} else {
				log.SetOutput(io.Discard)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = config.DefaultConfig()
			)

			return codegen.Generate(cfg)
		},
	}
)

func Execute() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.verbose, "verbose", "v", false, "shows the logs")
	cobra.CheckErr(rootCmd.Execute())
}
