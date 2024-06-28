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
	sqlDatabaser, sqlTabler, sqlColumner, sqlValuer, sqlScanner,
	binaryMarshaler, binaryUnmarshaler,
	textMarshaler, textUnmarshaler *types.Interface

	//go:embed sequel.go.tpl
	sqlBytes []byte
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, "database/sql...", "encoding")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		switch p.ID {
		case "database/sql/driver":
			goSqlValuer = p.Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
		case "database/sql":
			goSqlScanner = p.Types.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
		case "encoding":
			binaryMarshaler = p.Types.Scope().Lookup("BinaryMarshaler").Type().Underlying().(*types.Interface)
			binaryUnmarshaler = p.Types.Scope().Lookup("BinaryUnmarshaler").Type().Underlying().(*types.Interface)
			textMarshaler = p.Types.Scope().Lookup("TextMarshaler").Type().Underlying().(*types.Interface)
			textUnmarshaler = p.Types.Scope().Lookup("TextUnmarshaler").Type().Underlying().(*types.Interface)
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
	sqlDatabaser = pkg.Scope().Lookup("Databaser").Type().Underlying().(*types.Interface)
	sqlTabler = pkg.Scope().Lookup("Tabler").Type().Underlying().(*types.Interface)
	sqlColumner = pkg.Scope().Lookup("Columner").Type().Underlying().(*types.Interface)
	sqlValuer = pkg.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
	sqlScanner = pkg.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
}
