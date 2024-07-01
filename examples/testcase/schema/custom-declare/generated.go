package customdeclare

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
)

func (A) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{}
}
func (A) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO `a` (`name`) VALUES (?);", v.Values()
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("`name`", v.Name, func(val string) driver.Value { return string(val) })
}
