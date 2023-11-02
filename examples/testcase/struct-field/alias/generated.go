// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (AliasStruct) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS alias_struct (b VARCHAR(255) NOT NULL,Id BIGINT NOT NULL,header VARCHAR(255) NOT NULL,raw BLOB,text VARCHAR(255) NOT NULL,null_str VARCHAR(255) NOT NULL,created DATETIME NOT NULL,updated DATETIME NOT NULL,PRIMARY KEY (Id));"
}
func (AliasStruct) AlterTableStmt() string {
	return "ALTER TABLE alias_struct MODIFY b VARCHAR(255) NOT NULL,MODIFY Id BIGINT NOT NULL AFTER b,MODIFY header VARCHAR(255) NOT NULL AFTER Id,MODIFY raw BLOB AFTER header,MODIFY text VARCHAR(255) NOT NULL AFTER raw,MODIFY null_str VARCHAR(255) NOT NULL AFTER text,MODIFY created DATETIME NOT NULL AFTER null_str,MODIFY updated DATETIME NOT NULL AFTER created;"
}
func (AliasStruct) TableName() string {
	return "alias_struct"
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"}
}
func (v AliasStruct) IsAutoIncr() bool {
	return false
}
func (v AliasStruct) PK() (columnName string, pos int, value driver.Value) {
	return "Id", 1, int64(v.pk.ID)
}
func (v AliasStruct) Values() []any {
	return []any{float64(v.B), int64(v.pk.ID), string(v.Header), v.Raw, string(v.Text), (driver.Valuer)(v.NullStr), time.Time(v.model.Created), time.Time(v.model.Updated)}
}
func (v *AliasStruct) Addrs() []any {
	return []any{types.Float(&v.B), types.Integer(&v.pk.ID), types.String(&v.Header), &v.Raw, types.String(&v.Text), (sql.Scanner)(&v.NullStr), (*time.Time)(&v.model.Created), (*time.Time)(&v.model.Updated)}
}

func (B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS b (name VARCHAR(255) NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE b MODIFY name VARCHAR(255) NOT NULL;"
}
func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"name"}
}
func (v B) Values() []any {
	return []any{string(v.Name)}
}
func (v *B) Addrs() []any {
	return []any{types.String(&v.Name)}
}

func (C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS c (id BIGINT NOT NULL);"
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
func (v C) Values() []any {
	return []any{int64(v.ID)}
}
func (v *C) Addrs() []any {
	return []any{types.Integer(&v.ID)}
}
