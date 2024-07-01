package composite

import (
	"database/sql"
	"database/sql/driver"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Composite) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (Composite) TableName() string {
	return "`composite`"
}
func (Composite) HasPK() {}
func (v Composite) CompositeKey() ([]string, []int, []any) {
	return []string{"`col_1`", "`col_3`"}, []int{1, 3}, []any{string(v.Col1), (driver.Valuer)(v.Col3)}
}
func (Composite) ColumnNames() []string {
	return []string{"`flag`", "`col_1`", "`col_2`", "`col_3`"}
}
func (v Composite) Values() []any {
	return []any{bool(v.Flag), string(v.Col1), bool(v.Col2), (driver.Valuer)(v.Col3)}
}
func (v *Composite) Addrs() []any {
	return []any{types.Bool(&v.Flag), types.String(&v.Col1), types.Bool(&v.Col2), (sql.Scanner)(&v.Col3)}
}
func (Composite) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v Composite) InsertOneStmt() (string, []any) {
	return "INSERT INTO `composite` (`flag`,`col_1`,`col_2`,`col_3`) VALUES (?,?,?,?);", v.Values()
}
func (v Composite) FindOneByPKStmt() (string, []any) {
	return "SELECT `flag`,`col_1`,`col_2`,`col_3` FROM `composite` WHERE `col_1` = ? AND `col_3` = ? LIMIT 1;", []any{string(v.Col1), (driver.Valuer)(v.Col3)}
}
func (v Composite) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `composite` SET `flag` = ?,`col_2` = ? WHERE `col_1` = ? AND `col_3` = ? LIMIT 1;", []any{bool(v.Flag), bool(v.Col2), string(v.Col1), (driver.Valuer)(v.Col3)}
}
func (v Composite) GetFlag() sequel.ColumnValuer[bool] {
	return sequel.Column("`flag`", v.Flag, func(val bool) driver.Value { return bool(val) })
}
func (v Composite) GetCol1() sequel.ColumnValuer[string] {
	return sequel.Column("`col_1`", v.Col1, func(val string) driver.Value { return string(val) })
}
func (v Composite) GetCol2() sequel.ColumnValuer[bool] {
	return sequel.Column("`col_2`", v.Col2, func(val bool) driver.Value { return bool(val) })
}
func (v Composite) GetCol3() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("`col_3`", v.Col3, func(val uuid.UUID) driver.Value { return (driver.Valuer)(val) })
}
