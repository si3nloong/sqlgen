package cmd

import (
	"fmt"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version string",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(sequel.Version)
		},
	}
)
