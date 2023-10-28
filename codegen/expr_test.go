package codegen

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpr(t *testing.T) {
	pkg := new(Package)

	require.Equal(t, "string(v)", Expr(`string(%s)`).Format(pkg, "v"))
	require.Equal(t, `time.Time(v)`, Expr(`time.Time(%s)`).Format(pkg, "v"))
	require.Equal(t, `(*time.Time)(v)`, Expr(`(*time.Time)(%s)`).Format(pkg, "v"))
	require.Equal(t, `(driver.Valuer)(v)`, Expr(`(database/sql/driver.Valuer)(%s)`).Format(pkg, "v"))

	require.ElementsMatch(t, []*types.Package{
		types.NewPackage("time", "time"),
		types.NewPackage("database/sql/driver", "driver"),
	}, pkg.importPkgs)
}
