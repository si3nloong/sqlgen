package codegen

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	_ "embed"

	"golang.org/x/tools/go/packages"
)

var (
	// native package interface
	goSqlValuer, goSqlScanner,

	// sqlgen interface
	sqlDatabaser, sqlTabler, sqlColumner, sqlQueryColumner, sqlValuer, sqlScanner,
	binaryMarshaler, binaryUnmarshaler,
	textMarshaler, textUnmarshaler *types.Interface

	typeOfTime string

	//go:embed sequel.go.tpl
	sqlBytes []byte
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, "time", "sync", "database/sql...", "encoding")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		scope := p.Types.Scope()
		switch p.ID {
		case "database/sql/driver":
			goSqlValuer = scope.Lookup("Valuer").Type().Underlying().(*types.Interface)
		case "database/sql":
			goSqlScanner = scope.Lookup("Scanner").Type().Underlying().(*types.Interface)
		case "encoding":
			binaryMarshaler = scope.Lookup("BinaryMarshaler").Type().Underlying().(*types.Interface)
			binaryUnmarshaler = scope.Lookup("BinaryUnmarshaler").Type().Underlying().(*types.Interface)
			textMarshaler = scope.Lookup("TextMarshaler").Type().Underlying().(*types.Interface)
			textUnmarshaler = scope.Lookup("TextUnmarshaler").Type().Underlying().(*types.Interface)
		case "time":
			typeOfTime = scope.Lookup("Time").Type().(*types.Named).String()
		}
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", sqlBytes, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("sequel", fset, []*ast.File{f}, nil)
	if err != nil {
		panic(err)
	}
	scope := pkg.Scope()
	sqlDatabaser = scope.Lookup("Databaser").Type().Underlying().(*types.Interface)
	sqlTabler = scope.Lookup("Tabler").Type().Underlying().(*types.Interface)
	sqlColumner = scope.Lookup("Columner").Type().Underlying().(*types.Interface)
	sqlQueryColumner = scope.Lookup("SQLColumner").Type().Underlying().(*types.Interface)
	sqlValuer = scope.Lookup("Valuer").Type().Underlying().(*types.Interface)
	sqlScanner = scope.Lookup("Scanner").Type().Underlying().(*types.Interface)
}
