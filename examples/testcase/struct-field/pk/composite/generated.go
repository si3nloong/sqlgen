package composite

import (
	"database/sql"
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Composite) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `composite` (`flag` BOOL NOT NULL,`col_1` VARCHAR(255) NOT NULL,`col_2` VARCHAR(36),PRIMARY KEY (`col_1`,`col_2`));"
}
func (Composite) TableName() string {
	return "composite"
}
func (Composite) InsertOneStmt() string {
	return "INSERT INTO composite (flag,col_1,col_2) VALUES (?,?,?);"
}
func (Composite) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Composite) Columns() []string {
	return []string{"flag", "col_1", "col_2"}
}
func (Composite) HasPK() {}
func (v Composite) CompositeKey() ([]string, []int, []any) {
	return []string{"col_1", "col_2"}, []int{1, 2}, []any{string(v.Col1), (driver.Valuer)(v.Col2)}
}
func (Composite) FindByPKStmt() string {
	return "SELECT flag,col_1,col_2 FROM composite WHERE col_1 = ? AND col_2 = ? LIMIT 1;"
}
func (Composite) UpdateByPKStmt() string {
	return "UPDATE composite SET flag = ? WHERE col_1 = ? AND col_2 = ? LIMIT 1;"
}
func (v Composite) Values() []any {
	return []any{bool(v.Flag), string(v.Col1), (driver.Valuer)(v.Col2)}
}
func (v *Composite) Addrs() []any {
	return []any{types.Bool(&v.Flag), types.String(&v.Col1), (sql.Scanner)(&v.Col2)}
}
func (v Composite) GetFlag() sequel.ColumnValuer[bool] {
	return sequel.Column("flag", v.Flag, func(val bool) driver.Value { return bool(val) })
}
func (v Composite) GetCol1() sequel.ColumnValuer[string] {
	return sequel.Column("col_1", v.Col1, func(val string) driver.Value { return string(val) })
}
func (v Composite) GetCol2() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("col_2", v.Col2, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
