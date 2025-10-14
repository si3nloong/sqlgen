package main

import (
	"iter"
	"log"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
	"github.com/spf13/cobra"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate [path]",
		Short: "Print the version string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Args ->", args)
			path := args[0]
			matcher, err := codegen.PathResolver(path)
			if err != nil {
				return err
			}

			goPkg, tables, err := compiler.ParseDir(path, &compiler.Config{
				Matcher: matcher,
			})
			if err != nil {
				return err
			}

			_ = goPkg
			next, stop := iter.Pull2(tables)
			defer stop()

			for {
				t, err, ok := next()
				if err != nil {
					return err
				} else if !ok {
					return nil
				} else {
					log.Println(t)
				}
			}
		},
	}
)
