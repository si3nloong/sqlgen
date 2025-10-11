package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate create",
		Short: "Print the version string",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Args ->", args)
		},
	}
)
