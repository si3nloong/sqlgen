package customdeclare

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`name` VARCHAR(255) NOT NULL);"
}
func (A) AlterTableStmt() string {
	return "ALTER TABLE `a` MODIFY `name` VARCHAR(255) NOT NULL;"
}
func (A) InsertVarQuery() string {
	return "(?)"
}
func (v A) Values() []any {
	return []any{string(v.Name)}
}
func (v *A) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`name`", v.Name, func(vi string) driver.Value { return string(vi) })
}
