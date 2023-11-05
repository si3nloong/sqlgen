// Code generated by sqlgen, version v1.0.0-alpha. DO NOT EDIT.

package custom

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Address) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `address` (`line_1` VARCHAR(255) NOT NULL,`line_2` VARCHAR(255) NOT NULL,`city` VARCHAR(255) NOT NULL,`post_code` INTEGER UNSIGNED NOT NULL,`state_code` VARCHAR(255) NOT NULL,`country_code` VARCHAR(255) NOT NULL);"
}
func (Address) AlterTableStmt() string {
	return "ALTER TABLE `address` MODIFY `line_1` VARCHAR(255) NOT NULL,MODIFY `line_2` VARCHAR(255) NOT NULL AFTER `line_1`,MODIFY `city` VARCHAR(255) NOT NULL AFTER `line_2`,MODIFY `post_code` INTEGER UNSIGNED NOT NULL AFTER `city`,MODIFY `state_code` VARCHAR(255) NOT NULL AFTER `post_code`,MODIFY `country_code` VARCHAR(255) NOT NULL AFTER `state_code`;"
}
func (Address) TableName() string {
	return "`address`"
}
func (Address) Columns() []string {
	return []string{"`line_1`", "`line_2`", "`city`", "`post_code`", "`state_code`", "`country_code`"}
}
func (v Address) Values() []any {
	return []any{string(v.Line1), (driver.Valuer)(v.Line2), string(v.City), int64(v.PostCode), string(v.StateCode), string(v.CountryCode)}
}
func (v *Address) Addrs() []any {
	return []any{types.String(&v.Line1), (sql.Scanner)(&v.Line2), types.String(&v.City), types.Integer(&v.PostCode), types.String(&v.StateCode), types.String(&v.CountryCode)}
}

func (Customer) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `customer` (`id` BIGINT NOT NULL,`howOld` TINYINT UNSIGNED NOT NULL,`name` VARCHAR(255) NOT NULL,`address` JSON NOT NULL,`nicknames` JSON NOT NULL,`status` VARCHAR(255) NOT NULL,`join_at` DATETIME NOT NULL);"
}
func (Customer) AlterTableStmt() string {
	return "ALTER TABLE `customer` MODIFY `id` BIGINT NOT NULL,MODIFY `howOld` TINYINT UNSIGNED NOT NULL AFTER `id`,MODIFY `name` VARCHAR(255) NOT NULL AFTER `howOld`,MODIFY `address` JSON NOT NULL AFTER `name`,MODIFY `nicknames` JSON NOT NULL AFTER `address`,MODIFY `status` VARCHAR(255) NOT NULL AFTER `nicknames`,MODIFY `join_at` DATETIME NOT NULL AFTER `status`;"
}
func (Customer) TableName() string {
	return "`customer`"
}
func (Customer) Columns() []string {
	return []string{"`id`", "`howOld`", "`name`", "`address`", "`nicknames`", "`status`", "`join_at`"}
}
func (v Customer) Values() []any {
	return []any{int64(v.ID), int64(v.Age), (driver.Valuer)(v.Name), (driver.Valuer)(v.Address), encoding.MarshalStringList(v.Nicknames), string(v.Status), time.Time(v.JoinAt)}
}
func (v *Customer) Addrs() []any {
	return []any{types.Integer(&v.ID), types.Integer(&v.Age), (sql.Scanner)(&v.Name), &v.Address, types.StringList(&v.Nicknames), types.String(&v.Status), (*time.Time)(&v.JoinAt)}
}
