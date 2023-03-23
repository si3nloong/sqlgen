package cli

import (
	"github.com/spf13/cobra"
)

var (
	option = struct {
		watch bool
		force bool
	}{}

	rootCmd = &cobra.Command{
		Use:   "sqlgen",
		Short: "`sqlgen` is a SQL model generator.",
		// Long:  `A sql model generator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			codegen := newCodegen()
			if option.watch {
				watcher()
			}
			// codegen.run("./testdata/datatype/valuer-property/model.go")
			// codegen.run("./testdata/datatype/primitive-struct/model.go")
			// codegen.run("./testdata/datatype/custom-primitive-struct/model.go")
			// codegen.run("./testdata/datatype/empty-struct/model.go")
			// codegen.run("./testdata/datatype/array-property/customer.go")
			codegen.generate("./testdata/datatype/alias-property/common.go")
			// codegen.run("./testdata/datatype/pointer-property/model.go")
			return nil
		},
	}
)

func Execute() {
	// watcher
	rootCmd.Flags().BoolVarP(&option.watch, "watch", "w", false, "Watch the file changes and re-generate.")
	// force reload
	rootCmd.Flags().BoolVarP(&option.force, "force", "", false, "Force to re-generate.")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
