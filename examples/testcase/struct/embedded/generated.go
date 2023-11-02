// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package embedded

import (
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `b` (`id` BIGINT NOT NULL,`name` VARCHAR(255) NOT NULL,`z` TINYINT NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE `b` MODIFY `id` BIGINT NOT NULL,MODIFY `name` VARCHAR(255) NOT NULL AFTER `id`,MODIFY `z` TINYINT NOT NULL AFTER `name`;"
}
func (B) TableName() string {
	return "`b`"
}
func (B) Columns() []string {
	return []string{"`id`", "`name`", "`z`"}
}
func (v B) Values() []any {
	return []any{int64(v.a.ID), string(v.a.Name), bool(v.a.Z)}
}
func (v *B) Addrs() []any {
	return []any{types.Integer(&v.a.ID), types.String(&v.a.Name), types.Bool(&v.a.Z)}
}
