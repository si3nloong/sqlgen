// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package enum

import (
	"github.com/si3nloong/sqlgen/sql/types"
)

// Implements `sql.Valuer` interface.
func (Custom) Table() string {
	return "custom"
}

func (Custom) Columns() []string {
	return []string{"text", "e", "num"}
}

func (v Custom) Values() []any {
	return []any{string(v.Str), int64(v.Enum), int64(v.Num)}
}

func (v *Custom) Addrs() []any {
	return []any{types.String(&v.Str), types.Integer(&v.Enum), types.Integer(&v.Num)}
}
