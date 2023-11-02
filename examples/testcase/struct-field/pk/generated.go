// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package pk

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Car) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS car (id BIGINT NOT NULL,no VARCHAR(255) NOT NULL,color INTEGER NOT NULL,manuc_date DATETIME NOT NULL,PRIMARY KEY (id));"
}
func (Car) AlterTableStmt() string {
	return "ALTER TABLE car MODIFY id BIGINT NOT NULL,MODIFY no VARCHAR(255) NOT NULL AFTER id,MODIFY color INTEGER NOT NULL AFTER no,MODIFY manuc_date DATETIME NOT NULL AFTER color;"
}
func (Car) TableName() string {
	return "car"
}
func (Car) Columns() []string {
	return []string{"id", "no", "color", "manuc_date"}
}
func (v Car) IsAutoIncr() bool {
	return false
}
func (v Car) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, (driver.Valuer)(v.ID)
}
func (v Car) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.No), int64(v.Color), time.Time(v.ManucDate)}
}
func (v *Car) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No), types.Integer(&v.Color), (*time.Time)(&v.ManucDate)}
}

func (User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS user (id BIGINT NOT NULL,name VARCHAR(255) NOT NULL,age TINYINT UNSIGNED NOT NULL,email VARCHAR(255) NOT NULL,PRIMARY KEY (id));"
}
func (User) AlterTableStmt() string {
	return "ALTER TABLE user MODIFY id BIGINT NOT NULL,MODIFY name VARCHAR(255) NOT NULL AFTER id,MODIFY age TINYINT UNSIGNED NOT NULL AFTER name,MODIFY email VARCHAR(255) NOT NULL AFTER age;"
}
func (User) TableName() string {
	return "user"
}
func (User) Columns() []string {
	return []string{"id", "name", "age", "email"}
}
func (v User) IsAutoIncr() bool {
	return false
}
func (v User) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, int64(v.ID)
}
func (v User) Values() []any {
	return []any{int64(v.ID), string(v.Name), int64(v.Age), string(v.Email)}
}
func (v *User) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name), types.Integer(&v.Age), types.String(&v.Email)}
}

func (House) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS house (id INTEGER UNSIGNED NOT NULL,no VARCHAR(255) NOT NULL,PRIMARY KEY (id));"
}
func (House) AlterTableStmt() string {
	return "ALTER TABLE house MODIFY id INTEGER UNSIGNED NOT NULL,MODIFY no VARCHAR(255) NOT NULL AFTER id;"
}
func (House) TableName() string {
	return "house"
}
func (House) Columns() []string {
	return []string{"id", "no"}
}
func (v House) IsAutoIncr() bool {
	return false
}
func (v House) PK() (columnName string, pos int, value driver.Value) {
	return "id", 0, int64(v.ID)
}
func (v House) Values() []any {
	return []any{int64(v.ID), string(v.No)}
}
func (v *House) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No)}
}
