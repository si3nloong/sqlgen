package enum

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Custom) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`text`", Definition: "`text` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`e`", Definition: "`e` INTEGER NOT NULL DEFAULT 0"},
			{Name: "`num`", Definition: "`num` SMALLINT UNSIGNED NOT NULL DEFAULT 0"},
		},
	}
}
func (Custom) TableName() string {
	return "`custom`"
}
func (Custom) ColumnNames() []string {
	return []string{"`text`", "`e`", "`num`"}
}
func (v Custom) Values() []any {
	return []any{string(v.Str), int64(v.Enum), int64(v.Num)}
}
func (v *Custom) Addrs() []any {
	return []any{types.String(&v.Str), types.Integer(&v.Enum), types.Integer(&v.Num)}
}
func (Custom) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Custom) InsertOneStmt() (string, []any) {
	return "INSERT INTO `custom` (`text`,`e`,`num`) VALUES (?,?,?);", v.Values()
}
func (v Custom) GetStr() sequel.ColumnValuer[longText] {
	return sequel.Column("`text`", v.Str, func(val longText) driver.Value { return string(val) })
}
func (v Custom) GetEnum() sequel.ColumnValuer[Enum] {
	return sequel.Column("`e`", v.Enum, func(val Enum) driver.Value { return int64(val) })
}
func (v Custom) GetNum() sequel.ColumnValuer[uint16] {
	return sequel.Column("`num`", v.Num, func(val uint16) driver.Value { return int64(val) })
}
