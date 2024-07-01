package main

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
)

func (A) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (A) TableName() string {
	return "`a`"
}
func (A) ColumnNames() []string {
	return []string{"`t`"}
}
func (v A) Values() []any {
	return []any{time.Time(v.T)}
}
func (v *A) Addrs() []any {
	return []any{(*time.Time)(&v.T)}
}
func (A) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO `a` (`t`) VALUES (?);", v.Values()
}
func (v A) GetT() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`t`", v.T, func(val time.Time) driver.Value { return time.Time(val) })
}
