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
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = config.DefaultConfig()
				err error
			)

			// If user passing config file, then we load from it.
			if genOpts.config != "" {
				cfg, err = config.LoadConfigFrom(genOpts.config)
				if err != nil {
					return err
				}
			} else if len(args) > 0 {
				// If user pass the source, then we refer to it.
				cfg.Source = append(cfg.Source, args[0])
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
