package enum

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Custom) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`text` VARCHAR(255) NOT NULL,`e` INTEGER NOT NULL,`num` SMALLINT UNSIGNED NOT NULL);"
}
func (Custom) AlterTableStmt() string {
	return "ALTER TABLE `custom` MODIFY `text` VARCHAR(255) NOT NULL,MODIFY `e` INTEGER NOT NULL AFTER `text`,MODIFY `num` SMALLINT UNSIGNED NOT NULL AFTER `e`;"
}
func (Custom) TableName() string {
	return "`custom`"
}
func (v Custom) InsertOneStmt() string {
	return "INSERT INTO `custom` (`text`,`e`,`num`) VALUES (?,?,?);"
}
func (Custom) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Custom) Columns() []string {
	return []string{"`text`", "`e`", "`num`"}
}
func (v Custom) Values() []any {
	return []any{string(v.Str), int64(v.Enum), int64(v.Num)}
}
func (v *Custom) Addrs() []any {
	return []any{types.String(&v.Str), types.Integer(&v.Enum), types.Integer(&v.Num)}
}
func (v Custom) GetStr() sequel.ColumnValuer[longText] {
	return sequel.Column[longText]("`text`", v.Str, func(vi longText) driver.Value { return string(vi) })
}
func (v Custom) GetEnum() sequel.ColumnValuer[Enum] {
	return sequel.Column[Enum]("`e`", v.Enum, func(vi Enum) driver.Value { return int64(vi) })
}
func (v Custom) GetNum() sequel.ColumnValuer[uint16] {
	return sequel.Column[uint16]("`num`", v.Num, func(vi uint16) driver.Value { return int64(vi) })
}
