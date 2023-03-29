// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package primitive

import (
	"github.com/si3nloong/sqlgen/sql/types"
)

// Implements `sql.Valuer` interface.
func (Primitive) Table() string {
	return "primitive"
}

func (Primitive) Columns() []string {
	return []string{"str", "bytes", "bool", "int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "f_32", "f_64", "time"}
}

func (v Primitive) Values() []any {
	return []any{v.Str, string(v.Bytes), v.Bool, int64(v.Int), int64(v.Int8), int64(v.Int16), int64(v.Int32), v.Int64, int64(v.Uint), int64(v.Uint8), int64(v.Uint16), int64(v.Uint32), int64(v.Uint64), float64(v.F32), v.F64, v.Time}
}

func (v *Primitive) Addrs() []any {
	return []any{&v.Str, types.String(&v.Bytes), &v.Bool, types.Integer(&v.Int), types.Integer(&v.Int8), types.Integer(&v.Int16), types.Integer(&v.Int32), &v.Int64, types.Integer(&v.Uint), types.Integer(&v.Uint8), types.Integer(&v.Uint16), types.Integer(&v.Uint32), types.Integer(&v.Uint64), types.Float(&v.F32), &v.F64, &v.Time}
}
