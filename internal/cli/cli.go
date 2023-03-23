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
			// codegen.Generate("./testdata/datatype/valuer-property/model.go")
			// codegen.Generate("./testdata/datatype/primitive-struct/model.go")
			// codegen.Generate("./testdata/datatype/custom-primitive-struct/model.go")
			// codegen.Generate("./testdata/datatype/empty-struct/model.go")
			codegen.Generate("./testdata/datatype/array-property/customer.go")
			// codegen.Generate("./testdata/datatype/alias-property/common.go")
			// codegen.Generate("./testdata/datatype/pointer-property/model.go")
			return nil
		},
	}
)

func Execute() {
	// watcher
	rootCmd.Flags().BoolVarP(&option.watch, "watch", "w", false, "Watch the file changes and re-generate.")
	// force to reload
	rootCmd.Flags().BoolVarP(&option.force, "force", "", false, "Force to re-generate.")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
