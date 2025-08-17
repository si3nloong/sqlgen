package customdeclare

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
)

func (A) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (name) VALUES (?);", v.Values()
}
func (v A) NameValue() driver.Value {
	return v.Name
}
func (v A) ColumnName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
}
