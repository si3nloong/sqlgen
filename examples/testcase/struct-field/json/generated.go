package json

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (JSON) TableName() string {
	return "json"
}
func (JSON) Columns() []string {
	return []string{"num", "raw_bytes"} // 2
}
func (v JSON) Values() []any {
	return []any{
		v.Num.String(),     // 0 - num
		string(v.RawBytes), // 1 - raw_bytes
	}
}
func (v *JSON) Addrs() []any {
	return []any{
		Number(&v.Num), // 0 - num
		encoding.StringScanner[json.RawMessage](&v.RawBytes), // 1 - raw_bytes
	}
}
func (JSON) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v JSON) InsertOneStmt() (string, []any) {
	return "INSERT INTO json (num,raw_bytes) VALUES (?,?);", v.Values()
}
func (v JSON) NumValue() driver.Value {
	return v.Num.String()
}
func (v JSON) RawBytesValue() driver.Value {
	return string(v.RawBytes)
}
func (v JSON) GetNum() sequel.ColumnValuer[json.Number] {
	return sequel.Column("num", v.Num, func(val json.Number) driver.Value {
		return val.String()
	})
}
func (v JSON) GetRawBytes() sequel.ColumnValuer[json.RawMessage] {
	return sequel.Column("raw_bytes", v.RawBytes, func(val json.RawMessage) driver.Value {
		return string(val)
	})
}
