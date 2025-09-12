package inlinenested

import (
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{"nested"} // 1
}
func (v Model) Values() []any {
	return []any{
		encoding.JSONValue(v.Nested), // 0 - nested
	}
}
func (v *Model) Addrs() []any {
	return []any{
		encoding.JSONScanner(&v.Nested), // 0 - nested
	}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO `model` (`nested`) VALUES (?);", v.Values()
}
func (v Model) NestedValue() any {
	return encoding.JSONValue(v.Nested)
}

type ModelNested1 = struct {
	Time    time.Time
	Decimal decimal.Decimal
	Bool    bool
}

func (v Model) ColumnNested() sequel.ColumnValuer[ModelNested1] {
	return sequel.Column("nested", v.Nested, func(val ModelNested1) driver.Value {
		return encoding.JSONValue(val)
	})
}
