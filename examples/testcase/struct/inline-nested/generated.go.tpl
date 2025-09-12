package inlinenested

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (DeepNestedModel) TableName() string {
	return "deep_nested_model"
}
func (DeepNestedModel) Columns() []string {
	return []string{"nested"} // 1
}
func (v DeepNestedModel) Values() []any {
	return []any{
		encoding.JSONValue(v.Nested), // 0 - nested
	}
}
func (v *DeepNestedModel) Addrs() []any {
	return []any{
		encoding.JSONScanner(&v.Nested), // 0 - nested
	}
}
func (DeepNestedModel) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v DeepNestedModel) InsertOneStmt() (string, []any) {
	return "INSERT INTO `deep_nested_model` (`nested`) VALUES (?);", v.Values()
}
func (v DeepNestedModel) NestedValue() any {
	return encoding.JSONValue(v.Nested)
}

type DeepNestedModelNestedField = struct {
	Byte    sql.NullByte `sql:"byte" json:"byte"`
	Time    time.Time
	Decimal decimal.Decimal
	Bool    bool
}

func (v DeepNestedModel) ColumnNested() sequel.ColumnValuer[DeepNestedModelNestedField] {
	return sequel.Column("nested", v.Nested, func(val DeepNestedModelNestedField) driver.Value {
		return encoding.JSONValue(val)
	})
}

func (NestedModel) TableName() string {
	return "nested_model"
}
func (NestedModel) Columns() []string {
	return []string{"nested"} // 1
}
func (v NestedModel) Values() []any {
	return []any{
		encoding.JSONValue(v.Nested), // 0 - nested
	}
}
func (v *NestedModel) Addrs() []any {
	return []any{
		encoding.JSONScanner(&v.Nested), // 0 - nested
	}
}
func (NestedModel) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v NestedModel) InsertOneStmt() (string, []any) {
	return "INSERT INTO `nested_model` (`nested`) VALUES (?);", v.Values()
}
func (v NestedModel) NestedValue() any {
	return encoding.JSONValue(v.Nested)
}

type NestedModelNestedField = struct {
	Time    time.Time
	Decimal decimal.Decimal
	Bool    bool
}

func (v NestedModel) ColumnNested() sequel.ColumnValuer[NestedModelNestedField] {
	return sequel.Column("nested", v.Nested, func(val NestedModelNestedField) driver.Value {
		return encoding.JSONValue(val)
	})
}

func (NestedModelWithTag) TableName() string {
	return "nested_model_with_tag"
}
func (NestedModelWithTag) Columns() []string {
	return []string{"nested"} // 1
}
func (v NestedModelWithTag) Values() []any {
	return []any{
		encoding.JSONValue(v.Nested), // 0 - nested
	}
}
func (v *NestedModelWithTag) Addrs() []any {
	return []any{
		encoding.JSONScanner(&v.Nested), // 0 - nested
	}
}
func (NestedModelWithTag) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v NestedModelWithTag) InsertOneStmt() (string, []any) {
	return "INSERT INTO `nested_model_with_tag` (`nested`) VALUES (?);", v.Values()
}
func (v NestedModelWithTag) NestedValue() any {
	return encoding.JSONValue(v.Nested)
}

type NestedModelWithTagNestedField = struct {
	Time    time.Time
	Decimal decimal.Decimal
	Bool    bool `json:"bool"`
	Str     Str  `json:"str"`
}

func (v NestedModelWithTag) ColumnNested() sequel.ColumnValuer[NestedModelWithTagNestedField] {
	return sequel.Column("nested", v.Nested, func(val NestedModelWithTagNestedField) driver.Value {
		return encoding.JSONValue(val)
	})
}
