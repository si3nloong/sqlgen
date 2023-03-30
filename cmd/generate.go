package cmd

import (
	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/spf13/cobra"
)

var (
	genOpts struct {
		config string
		watch  bool
		force  bool
	}

	genCmd = &cobra.Command{
		Use:   "generate [source]",
		Short: "Generate struct functions for target models.",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = config.DefaultConfig()
				err error
			)
			if genOpts.config != "" {
				cfg, err = config.LoadConfigFrom(genOpts.config)
				if err != nil {
					return err
				}
			}
			if len(args) > 0 {
				cfg.SrcDir = args[0]
			}
			return codegen.Generate(cfg)
		},
	}
)

func init() {
	genCmd.Flags().StringVarP(&genOpts.config, "file", "f", "", "Config file path.")
	genCmd.Flags().BoolVarP(&genOpts.watch, "watch", "w", false, "Watch the file changes and re-generate.")
	genCmd.Flags().BoolVarP(&genOpts.force, "force", "", false, "Force to re-generate.")
}
