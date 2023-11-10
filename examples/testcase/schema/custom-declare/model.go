package customdeclare

type A struct {
	Name string
}

// Codegen will not override the custom declaration
func (A) TableName() string {
	return "mytable"
}

// Codegen will not override the custom declaration
func (A) Columns() []string {
	return []string{`a`, "b", "b"}
}
