package codegen

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	sqlValuer, sqlScanner,
	sqlTabler, sqlColumner, sqlRower *types.Interface
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, "database/sql...", "github.com/si3nloong/sqlgen/sequel")
	if err != nil {
		panic(err)
	}

	sqlValuer = pkgs[0].Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
	sqlScanner = pkgs[1].Types.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
	sqlTabler = pkgs[2].Types.Scope().Lookup("Tabler").Type().Underlying().(*types.Interface)
	sqlColumner = pkgs[2].Types.Scope().Lookup("Columner").Type().Underlying().(*types.Interface)
	sqlRower = pkgs[2].Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
}
