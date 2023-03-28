package cmd

import (
	"log"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/spf13/cobra"
)

var (
	options = struct {
		watch bool
		force bool
	}{}

	rootCmd = &cobra.Command{
		Use:   "sqlgen",
		Short: "`sqlgen` is a SQL model generator.",
		// Long:  `A sql model generator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	cmd.AddCommand(initCommand())
			// }
			log.Println("args ->", args)
			if options.watch {
				watcher()
			}

			return codegen.Generate(config.DefaultConfig())
			// log.Println(args)
			// codegen := newCodegen()

			// codegen.Generate("./testdata/datatype/array-property/customer.go")
			// codegen.Generate("./testdata/schema/model.go")
			// files, err := filepath.Glob(args[0])
			// return nil
		},
	}
)

func Execute() {
	rootCmd.AddCommand(initCommand())
	// watcher
	rootCmd.Flags().BoolVarP(&options.watch, "watch", "w", false, "Watch the file changes and re-generate.")
	// force to reload
	rootCmd.Flags().BoolVarP(&options.force, "force", "", false, "Force to re-generate.")
	if err := rootCmd.Execute(); err != nil {
		// panic(err)
		log.Fatal(err)
	}
}
