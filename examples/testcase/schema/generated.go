package schema

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v A) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `Apple` (`id` VARCHAR(255) NOT NULL,`text` VARCHAR(255) NOT NULL,`created_at` DATETIME NOT NULL);"
}
func (A) TableName() string {
	return "Apple"
}
func (A) InsertOneStmt() string {
	return "INSERT INTO Apple (id,text,created_at) VALUES (?,?,?);"
}
func (A) InsertVarQuery() string {
	return "(?,?,?)"
}
func (A) Columns() []string {
	return []string{"id", "text", "created_at"}
}
func (v A) Values() []any {
	return []any{string(v.ID), string(v.Text), time.Time(v.CreatedAt)}
}
func (v *A) Addrs() []any {
	return []any{types.String(&v.ID), types.String(&v.Text), (*time.Time)(&v.CreatedAt)}
}
func (v A) GetID() sequel.ColumnValuer[string] {
	return sequel.Column("id", v.ID, func(val string) driver.Value { return string(val) })
}
func (v A) GetText() sequel.ColumnValuer[LongText] {
	return sequel.Column("text", v.Text, func(val LongText) driver.Value { return string(val) })
}
func (v A) GetCreatedAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created_at", v.CreatedAt, func(val time.Time) driver.Value { return time.Time(val) })
}

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `b` (`id` VARCHAR(255) NOT NULL,`created_at` DATETIME NOT NULL);"
}
func (B) TableName() string {
	return "b"
}
func (B) InsertOneStmt() string {
	return "INSERT INTO b (id,created_at) VALUES (?,?);"
}
func (B) InsertVarQuery() string {
	return "(?,?)"
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
func (v B) GetID() sequel.ColumnValuer[string] {
	return sequel.Column("id", v.ID, func(val string) driver.Value { return string(val) })
}
func (v B) GetCreatedAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created_at", v.CreatedAt, func(val time.Time) driver.Value { return time.Time(val) })
}

func (v C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `c` (`id` BIGINT NOT NULL,PRIMARY KEY (`id`));"
}
func (C) TableName() string {
	return "c"
}
func (C) InsertOneStmt() string {
	return "INSERT INTO c (id) VALUES (?);"
}
func (C) InsertVarQuery() string {
	return "(?)"
}
func (C) Columns() []string {
	return []string{"id"}
}
func (C) HasPK() {}
func (v C) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{int64(v.ID)}
}
func (v C) Values() []any {
	return []any{int64(v.ID)}
}
func (v *C) Addrs() []any {
	return []any{types.Integer(&v.ID)}
}
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}

func (v D) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `d` (`id` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (D) TableName() string {
	return "d"
}
func (D) InsertOneStmt() string {
	return "INSERT INTO d (id) VALUES (?);"
}
func (D) InsertVarQuery() string {
	return "(?)"
}
func (D) Columns() []string {
	return []string{"id"}
}
func (D) HasPK() {}
func (v D) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{(driver.Valuer)(v.ID)}
}
func (v D) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *D) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (v D) GetID() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("id", v.ID, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
