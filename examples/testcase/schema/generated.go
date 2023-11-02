// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package schema

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS model ();"
}
func (Model) AlterTableStmt() string {
	return "ALTER TABLE model ;"
}
func (Model) TableName() string {
	return "model"
}
func (Model) Columns() []string {
	return []string{}
}
func (v Model) Values() []any {
	return []any{}
}
func (v *Model) Addrs() []any {
	return []any{}
}

func (A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS a (id VARCHAR(255) NOT NULL,created_at DATETIME NOT NULL);"
}
func (A) AlterTableStmt() string {
	return "ALTER TABLE a MODIFY id VARCHAR(255) NOT NULL,MODIFY created_at DATETIME NOT NULL AFTER id;"
}
func (A) TableName() string {
	return "a"
}
func (A) Columns() []string {
	return []string{"id", "created_at"}
}
func (v A) Values() []any {
	return []any{string(v.ID), time.Time(v.CreatedAt)}
}
func (v *A) Addrs() []any {
	return []any{types.String(&v.ID), (*time.Time)(&v.CreatedAt)}
}

func (B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS b (id VARCHAR(255) NOT NULL,created_at DATETIME NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE b MODIFY id VARCHAR(255) NOT NULL,MODIFY created_at DATETIME NOT NULL AFTER id;"
}
func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "created_at"}
}
func (v B) Values() []any {
	return []any{string(v.ID), time.Time(v.CreatedAt)}
}
func (v *B) Addrs() []any {
	return []any{types.String(&v.ID), (*time.Time)(&v.CreatedAt)}
}

func (C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS c (id BIGINT NOT NULL,PRIMARY KEY (id));"
}
func (C) AlterTableStmt() string {
	return "ALTER TABLE c MODIFY id BIGINT NOT NULL;"
}
func (C) TableName() string {
	return "c"
}
func (C) Columns() []string {
	return []string{"id"}
}
func (v C) IsAutoIncr() bool {
	return false
}
func (v C) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, int64(v.ID)
}
func (v C) Values() []any {
	return []any{int64(v.ID)}
}
func (v *C) Addrs() []any {
	return []any{types.Integer(&v.ID)}
}

func (D) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS d (id VARCHAR(255) NOT NULL,PRIMARY KEY (id));"
}
func (D) AlterTableStmt() string {
	return "ALTER TABLE d MODIFY id VARCHAR(255) NOT NULL;"
}
func (D) TableName() string {
	return "d"
}
func (D) Columns() []string {
	return []string{"id"}
}
func (v D) IsAutoIncr() bool {
	return false
}
func (v D) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, (driver.Valuer)(v.ID)
}
func (v D) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *D) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
