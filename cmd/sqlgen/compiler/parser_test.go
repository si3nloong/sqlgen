package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestFindFirstConstValueInPackageFileOrder(t *testing.T) {
	sources := map[string]string{
		"model.go": `package example
type longText string
const otherStr longText = "model"
`,
		"file2.go": `package example
//line zzz.go:100
const otherStr2 longText = "first"
//line aaa.go:1
const alphabeticallyFirst longText = "second"
`,
	}
	for _, order := range [][]string{
		{"model.go", "file2.go"},
		{"file2.go", "model.go"},
	} {
		t.Run(order[0], func(t *testing.T) {
			fset := token.NewFileSet()
			var files []*ast.File
			for _, name := range order {
				file, err := parser.ParseFile(fset, name, sources[name], 0)
				if err != nil {
					t.Fatal(err)
				}
				files = append(files, file)
			}
			var config types.Config
			typed, err := config.Check("example", fset, files, nil)
			if err != nil {
				t.Fatal(err)
			}
			pkg := &packages.Package{Types: typed, Fset: fset}
			got, ok := findFirstConstValueInPackage(pkg, "longText")
			if !ok {
				t.Fatal("no constant found")
			}
			if got.Name() != "otherStr2" {
				t.Fatalf("got %s, want otherStr2", got.Name())
			}
			if _, ok := findFirstConstValueInPackage(pkg, "missing"); ok {
				t.Fatal("found constant for missing type")
			}
		})
	}
}
