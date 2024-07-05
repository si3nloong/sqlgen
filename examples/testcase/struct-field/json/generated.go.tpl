package json

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (JSON) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`num`", Definition: "`num` VARCHAR(20) NOT NULL"},
			{Name: "`raw_bytes`", Definition: "`raw_bytes` JSON NOT NULL"},
		},
	}
}
func (JSON) TableName() string {
	return "`json`"
}
func (JSON) ColumnNames() []string {
	return []string{"`num`", "`raw_bytes`"}
}
func (v JSON) Values() []any {
	return []any{v.Num.String(), string(v.RawBytes)}
}
func (v *JSON) Addrs() []any {
	return []any{Number(&v.Num), types.String(&v.RawBytes)}
}
func (JSON) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v JSON) InsertOneStmt() (string, []any) {
	return "INSERT INTO `json` (`num`,`raw_bytes`) VALUES (?,?);", v.Values()
}
func (v JSON) GetNum() sequel.ColumnValuer[json.Number] {
	return sequel.Column("`num`", v.Num, func(val json.Number) driver.Value { return val.String() })
}
func (v JSON) GetRawBytes() sequel.ColumnValuer[json.RawMessage] {
	return sequel.Column("`raw_bytes`", v.RawBytes, func(val json.RawMessage) driver.Value { return string(val) })
}
