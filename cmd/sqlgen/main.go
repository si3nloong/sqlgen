// Typically this process would be run using go generate, like this:
//
//	//go:generate sqlgen
package main

import (
	"log"
)

func main() {
	log.SetPrefix("sqlgen: ")
	Execute()
}
