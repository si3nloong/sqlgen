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
	sqlValuer, sqlScanner,
	sqlTabler, sqlColumner,
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
			sqlValuer = p.Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
		case "database/sql":
			sqlScanner = p.Types.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
		case "encoding":
			binaryMarshaler = p.Types.Scope().Lookup("BinaryMarshaler").Type().Underlying().(*types.Interface)
			binaryUnmarshaler = p.Types.Scope().Lookup("BinaryUnmarshaler").Type().Underlying().(*types.Interface)
			textMarshaler = p.Types.Scope().Lookup("TextMarshaler").Type().Underlying().(*types.Interface)
			textUnmarshaler = p.Types.Scope().Lookup("TextUnmarshaler").Type().Underlying().(*types.Interface)
			// case "github.com/si3nloong/sqlgen/sequel":
			// 	sqlTabler = p.Types.Scope().Lookup("Tabler").Type().Underlying().(*types.Interface)
			// 	sqlColumner = p.Types.Scope().Lookup("Columner").Type().Underlying().(*types.Interface)
			// sqlRower = p.Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
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
		sqlTabler = pkg.Scope().Lookup("Tabler").Type().Underlying().(*types.Interface)
		sqlColumner = pkg.Scope().Lookup("Columner").Type().Underlying().(*types.Interface)
	}
}
