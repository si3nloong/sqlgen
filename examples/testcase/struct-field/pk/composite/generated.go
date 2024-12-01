package composite

import (
	"database/sql/driver"
)

func (Composite) TableName() string {
	return "composite"
}
func (Composite) HasPK() {}
func (v Composite) CompositeKey() ([]string, []int, []any) {
	return []string{"col_1", "col_3"}, []int{1, 3}, []any{v.Col1, v.Col3}
}
func (Composite) Columns() []string {
	return []string{"flag", "col_1", "col_2", "col_3"} // 4
}
func (v Composite) Values() []any {
	return []any{
		v.Flag, // 0 - flag
		v.Col1, // 1 - col_1
		v.Col2, // 2 - col_2
		v.Col3, // 3 - col_3
	}
}
func (v *Composite) Addrs() []any {
	return []any{
		&v.Flag, // 0 - flag
		&v.Col1, // 1 - col_1
		&v.Col2, // 2 - col_2
		&v.Col3, // 3 - col_3
	}
}
func (Composite) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v Composite) InsertOneStmt() (string, []any) {
	return "INSERT INTO composite (flag,col_1,col_2,col_3) VALUES (?,?,?,?);", v.Values()
}
func (v Composite) FindOneByPKStmt() (string, []any) {
	return "SELECT flag,col_1,col_2,col_3 FROM composite WHERE (col_1,col_3) = (?,?) LIMIT 1;", []any{v.Col1, v.Col3}
}
func (v Composite) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE composite SET flag = ?,col_2 = ? WHERE (col_1,col_3) = (?,?);", []any{v.Flag, v.Col2, v.Col1, v.Col3}
}
func (v Composite) GetFlag() driver.Value {
	return v.Flag
}
func (v Composite) GetCol1() driver.Value {
	return v.Col1
}
func (v Composite) GetCol2() driver.Value {
	return v.Col2
}
func (v Composite) GetCol3() driver.Value {
	return v.Col3
}
