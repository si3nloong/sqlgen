package customdeclare

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v A) Values() []any {
	return []any{string(v.Name)}
}
func (v *A) Addrs() []any {
	return []any{types.String(&v.Name)}
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
