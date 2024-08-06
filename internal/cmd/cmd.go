package cmd

import (
	"io"
	"log"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/spf13/cobra"
)

var (
	rootOpts struct {
		config  string
		verbose bool
	}

	rootCmd = &cobra.Command{
		Use:          "sqlgen",
		Short:        "🚀 Transform your struct to Go code!!!",
		Long:         ``,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if rootOpts.verbose {
				log.SetFlags(0)
			} else {
				log.SetOutput(io.Discard)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = codegen.DefaultConfig()
				err error
			)

			// If user passing config file, then we load from it.
			if rootOpts.config != "" {
				cfg, err = codegen.LoadConfigFrom(rootOpts.config)
				if err != nil {
					return err
				}
			}
			return codegen.Generate(cfg)
		},
	}
)

func Execute() {
	log.SetPrefix("sqlgen:")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVarP(&rootOpts.config, "config", "c", "", "config file")
	// rootCmd.Flags().BoolVarP(&rootOpts.watch, "watch", "w", false, "watch the file changes and re-generate.")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.verbose, "verbose", "v", false, "shows the logs")
	cobra.CheckErr(rootCmd.Execute())
}
