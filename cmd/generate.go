package cmd

import (
	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/spf13/cobra"
)

var (
	genOpts struct {
		watch bool
		force bool
	}

	genCmd = &cobra.Command{
		Use:   "generate [source]",
		Short: "Generate boilerplate code based on go struct",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = config.DefaultConfig()
			)

			// If user pass the source, then we refer to it.
			cfg.Source = []string{args[0]}

			return codegen.Generate(cfg)
		},
	}
)

func init() {
	// genCmd.Flags().BoolVarP(&genOpts.watch, "watch", "w", false, "watch the file changes and re-generate.")
	genCmd.Flags().BoolVarP(&genOpts.force, "force", "f", false, "force to execute")
}
