package cmd

import (
	"github.com/si3nloong/sqlgen/codegen/config"
	"io"
	"log"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/spf13/cobra"
)

var (
	options = struct {
		watch   bool
		force   bool
		verbose bool
	}{}

	rootCmd = &cobra.Command{
		Use:   "sqlgen",
		Short: "`sqlgen` is a SQL model generator.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.verbose {
				log.SetFlags(0)
			} else {
				log.SetOutput(io.Discard)
			}

			return codegen.Generate(config.DefaultConfig())
		},
	}
)

func Execute() {
	rootCmd.AddCommand(initCommand())
	rootCmd.Flags().BoolVarP(&options.verbose, "verbose", "v", false, "Shows the log details.")
	// watcher
	rootCmd.Flags().BoolVarP(&options.watch, "watch", "w", false, "Watch the file changes and re-generate.")
	// force to regenerate
	rootCmd.Flags().BoolVarP(&options.force, "force", "", false, "Force to re-generate.")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
