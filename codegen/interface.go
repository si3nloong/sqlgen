package codegen

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	sqlValuer, sqlScanner,
	sqlTabler, sqlColumner, sqlRower,
	binaryMarshaler, binaryUnmarshaler,
	textMarshaler, textUnmarshaler *types.Interface
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, "database/sql...", "github.com/si3nloong/sqlgen/sequel", "encoding")
	if err != nil {
		panic(err)
	}

	sqlValuer = pkgs[0].Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
	sqlScanner = pkgs[1].Types.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
	sqlTabler = pkgs[2].Types.Scope().Lookup("Tabler").Type().Underlying().(*types.Interface)
	sqlColumner = pkgs[2].Types.Scope().Lookup("Columner").Type().Underlying().(*types.Interface)
	sqlRower = pkgs[2].Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
	binaryMarshaler = pkgs[3].Types.Scope().Lookup("BinaryMarshaler").Type().Underlying().(*types.Interface)
	binaryUnmarshaler = pkgs[3].Types.Scope().Lookup("BinaryUnmarshaler").Type().Underlying().(*types.Interface)
	textMarshaler = pkgs[3].Types.Scope().Lookup("TextMarshaler").Type().Underlying().(*types.Interface)
	textUnmarshaler = pkgs[3].Types.Scope().Lookup("TextUnmarshaler").Type().Underlying().(*types.Interface)
}
