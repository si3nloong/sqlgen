package aliasstruct

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`date` DATE NOT NULL,`time` VARCHAR(255) NOT NULL);"
}
func (A) InsertVarQuery() string {
	return "(?,?)"
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
func (v A) GetDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("date", v.Date, func(val civil.Date) driver.Value { return types.TextMarshaler(val) })
}
func (v A) GetTime() sequel.ColumnValuer[civil.Time] {
	return sequel.Column("time", v.Time, func(val civil.Time) driver.Value { return types.TextMarshaler(val) })
}

func (v C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `c` (`string` VARCHAR(255) NOT NULL,`valid` BOOL NOT NULL);"
}
func (C) TableName() string {
	return "c"
}
func (C) InsertOneStmt() string {
	return "INSERT INTO c (string,valid) VALUES (?,?);"
}
func (C) InsertVarQuery() string {
	return "(?,?)"
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
func (v C) GetString() sequel.ColumnValuer[string] {
	return sequel.Column("string", v.String, func(val string) driver.Value { return string(val) })
}
func (v C) GetValid() sequel.ColumnValuer[bool] {
	return sequel.Column("valid", v.Valid, func(val bool) driver.Value { return bool(val) })
}
