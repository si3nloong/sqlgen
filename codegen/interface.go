package codegen

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	sqlValuer, sqlScanner *types.Interface
)

func init() {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
	}, "database/sql...")
	if err != nil {
		panic(err)
	}

	sqlValuer = pkgs[0].Types.Scope().Lookup("Valuer").Type().Underlying().(*types.Interface)
	sqlScanner = pkgs[1].Types.Scope().Lookup("Scanner").Type().Underlying().(*types.Interface)
}
