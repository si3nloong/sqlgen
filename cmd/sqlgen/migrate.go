package main

import (
	"fmt"
	"iter"
	"os"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate [src] [outDir]",
		Short: "Print the version string",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			outDir := args[1]
			stat, err := os.Stat(outDir)
			if err != nil && !os.IsNotExist(err) {
				return err
			} else if stat != nil && !stat.IsDir() {
				return fmt.Errorf(`sqlgen: "outDir" must be a directory`)
			}

			if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
				return err
			}

			cfg := codegen.DefaultConfig()
			cfg.Source = []string{source}

			return codegen.Walk(cfg, func(gen *codegen.Generator, pkg *packages.Package, tables iter.Seq2[*compiler.Table, error]) error {
				return gen.GenerateMigrations(outDir, pkg, tables)
			})
		},
	}
)
