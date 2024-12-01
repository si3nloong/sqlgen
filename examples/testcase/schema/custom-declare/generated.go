package customdeclare

import (
	"database/sql/driver"
)

func (A) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (name) VALUES (?);", v.Values()
}
func (v A) GetName() driver.Value {
	return v.Name
}
