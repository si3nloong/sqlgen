package aliasstruct

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (A) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "date", Definition: "date DATE NOT NULL"},
			{Name: "time", Definition: "time TIME NOT NULL"},
		},
	}
}
func (A) Columns() []string {
	return []string{"date", "time"}
}
func (v A) Values() []any {
	return []any{types.TextMarshaler(v.Date), types.TextMarshaler(v.Time)}
}
func (v *A) Addrs() []any {
	return []any{types.Date(&v.Date), types.TextUnmarshaler(&v.Time)}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO a (date,time) VALUES (?,?);", v.Values()
}
func (v A) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.Date, func(val civil.Date) driver.Value { return types.TextMarshaler(val) })
}
func (v A) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.Time, func(val civil.Time) driver.Value { return types.TextMarshaler(val) })
}

func (C) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "string", Definition: "string VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "valid", Definition: "valid BOOL NOT NULL DEFAULT false"},
		},
	}
}
func (C) TableName() string {
	return "c"
}
func (C) Columns() []string {
	return []string{"string", "valid"}
}
func (v C) Values() []any {
	return []any{string(v.String), bool(v.Valid)}
}
func (v *C) Addrs() []any {
	return []any{types.String(&v.String), types.Bool(&v.Valid)}
}
func (C) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v C) InsertOneStmt() (string, []any) {
	return "INSERT INTO c (string,valid) VALUES (?,?);", v.Values()
}
func (v C) GetString() sequel.ColumnValuer[string] {
	return sequel.Column("string", v.String, func(val string) driver.Value { return string(val) })
}
func (v C) GetValid() sequel.ColumnValuer[bool] {
	return sequel.Column("valid", v.Valid, func(val bool) driver.Value { return bool(val) })
}
