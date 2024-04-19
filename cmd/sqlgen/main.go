// Typically this process would be run using go generate, like this:
//
//	//go:generate sqlgen
package main

import (
	"log"

	"github.com/si3nloong/sqlgen/internal/cmd"
)

func main() {
	log.SetPrefix("sqlgen: ")
	cmd.Execute()
}
