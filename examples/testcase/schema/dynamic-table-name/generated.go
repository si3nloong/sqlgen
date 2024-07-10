package tabler

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Model) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`name`", Definition: "`name` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (Model) Columns() []string {
	return []string{"`name`"}
}
func (v Model) Values() []any {
	return []any{string(v.Name)}
}
func (v *Model) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (Model) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v Model) InsertOneStmt() (string, []any) {
	return "INSERT INTO `model` (`name`) VALUES (?);", v.Values()
}
func (v Model) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("`name`", v.Name, func(val string) driver.Value { return string(val) })
}

func (A) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT NOT NULL"},
			{Name: "`name`", Definition: "`name` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (A) HasPK() {}
func (v A) PK() (string, int, any) {
	return "`id`", 0, int64(v.ID)
}
func (A) Columns() []string {
	return []string{"`id`", "`name`"}
}
func (v A) Values() []any {
	return []any{int64(v.ID), string(v.Name)}
}
func (v *A) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name)}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO `a` (`id`,`name`) VALUES (?,?);", v.Values()
}
func (v A) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`name` FROM `a` WHERE `id` = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v A) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `a` SET `name` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.Name), int64(v.ID)}
}
func (v A) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("`id`", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v A) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("`name`", v.Name, func(val string) driver.Value { return string(val) })
}
