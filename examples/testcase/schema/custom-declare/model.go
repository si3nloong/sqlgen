package customdeclare

import "github.com/si3nloong/sqlgen/sequel/types"

type A struct {
	Name string
}

// Codegen will not override the custom declaration
func (A) TableName() string {
	return "mytable"
}

// Codegen will not override the custom declaration even if it has error
func (A) Columns() []string {
	return []string{`a`, "b", "b"}
}

// Codegen will not override the custom declaration even if it has error
func (v A) Values() []any {
	return []any{string(v.Name)}
}

// Codegen will not override the custom declaration even if it has error
func (v *A) Addrs() []any {
	return []any{types.String(&v.Name)}
}
