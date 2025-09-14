package codegen

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpr(t *testing.T) {
	pkg := NewPackage("", "")

	require.Equal(t, "string(v)", Expr(`string({{goPath}})`).Format(pkg, ExprParams{GoPath: "v"}))
	require.Equal(t, `time.Time(v)`, Expr(`time.Time({{goPath}})`).Format(pkg, ExprParams{GoPath: "v"}))
	require.Equal(t, `(*time.Time)(&v)`, Expr(`(*time.Time)({{addrOfGoPath}})`).Format(pkg, ExprParams{GoPath: "v"}))
	require.Equal(t, `(driver.Valuer)(v)`, Expr(`(database/sql/driver.Valuer)({{goPath}})`).Format(pkg, ExprParams{GoPath: "v"}))

	require.ElementsMatch(t, []*types.Package{
		types.NewPackage("time", "time"),
		types.NewPackage("database/sql/driver", "driver"),
	}, pkg.imports)
}
