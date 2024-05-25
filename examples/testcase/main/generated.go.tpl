package main

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
)

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`t` DATETIME NOT NULL);"
}
func (v A) AlterTableStmt() string {
	return "ALTER TABLE " + v.TableName() + " (MODIFY `t` DATETIME NOT NULL);"
}
func (A) TableName() string {
	return "a"
}
func (v A) InsertOneStmt() string {
	return "INSERT INTO a (t) VALUES (?);"
}
func (A) InsertVarQuery() string {
	return "(?)"
}
func (A) Columns() []string {
	return []string{"t"}
}
func (v A) Values() []any {
	return []any{time.Time(v.T)}
}
func (v *A) Addrs() []any {
	return []any{(*time.Time)(&v.T)}
}
func (v A) GetT() sequel.ColumnValuer[time.Time] {
	return sequel.Column("t", v.T, func(vi time.Time) driver.Value { return time.Time(vi) })
}
