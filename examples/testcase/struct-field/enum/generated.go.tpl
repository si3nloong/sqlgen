package enum

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Custom) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `custom` (`text` VARCHAR(255) NOT NULL,`e` INTEGER NOT NULL,`num` SMALLINT UNSIGNED NOT NULL);"
}
func (Custom) TableName() string {
	return "custom"
}
func (Custom) InsertOneStmt() string {
	return "INSERT INTO custom (text,e,num) VALUES (?,?,?);"
}
func (Custom) InsertVarQuery() string {
	return "(?,?,?)"
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
func (v Custom) GetStr() sequel.ColumnValuer[longText] {
	return sequel.Column("text", v.Str, func(val longText) driver.Value { return string(val) })
}
func (v Custom) GetEnum() sequel.ColumnValuer[Enum] {
	return sequel.Column("e", v.Enum, func(val Enum) driver.Value { return int64(val) })
}
func (v Custom) GetNum() sequel.ColumnValuer[uint16] {
	return sequel.Column("num", v.Num, func(val uint16) driver.Value { return int64(val) })
}
