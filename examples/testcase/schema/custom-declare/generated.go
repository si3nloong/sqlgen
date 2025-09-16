package customdeclare

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (A) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO " + v.TableName() + " (`name`) VALUES (?);", v.Values()
}
func (v A) NameValue() any {
	return v.Name
}
func (v A) ColumnName() sequel.ColumnClause {
	return sequel.BasicColumn("name", v.Name)
}
