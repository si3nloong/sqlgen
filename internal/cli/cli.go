package cli

import (
	"github.com/spf13/cobra"
)

var watch bool
var rootCmd = &cobra.Command{
	Use:   "sqlgen",
	Short: "`sqlgen` is a SQL model generator.",
	// Long:  `A sql model generator.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		codegen := newCodegen()
		if watch {
			watcher()
		}
		// codegen.run("./testdata/datatype/valuer-property/model.go")
		// codegen.run("./testdata/datatype/primitive-struct/model.go")
		// codegen.run("./testdata/datatype/custom-primitive-struct/model.go")
		// codegen.run("./testdata/datatype/empty-struct/model.go")
		// codegen.run("./testdata/datatype/array-property/customer.go")
		codegen.run("./testdata/datatype/alias-property/common.go")
		// codegen.run("./testdata/datatype/pointer-property/model.go")
		return nil
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&watch, "watch", "w", false, "")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
