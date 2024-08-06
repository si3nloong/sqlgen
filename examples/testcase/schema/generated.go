package schema

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (A) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "id", Definition: "id VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "text", Definition: "text VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "created_at", Definition: "created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"},
		},
	}
}
func (A) TableName() string {
	return "Apple"
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
func (A) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO Apple (id,text,created_at) VALUES (?,?,?);", v.Values()
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

func (B) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "id", Definition: "id VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "created_at", Definition: "created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"},
		},
	}
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
func (B) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,created_at) VALUES (?,?);", v.Values()
}
func (v B) GetID() sequel.ColumnValuer[string] {
	return sequel.Column("id", v.ID, func(val string) driver.Value { return string(val) })
}
func (v B) GetCreatedAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created_at", v.CreatedAt, func(val time.Time) driver.Value { return time.Time(val) })
}

func (C) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"id"},
			Definition: "PRIMARY KEY (id)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "id", Definition: "id BIGINT NOT NULL"},
		},
	}
}
func (C) TableName() string {
	return "c"
}
func (C) HasPK() {}
func (v C) PK() (string, int, any) {
	return "id", 0, int64(v.ID)
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
func (C) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v C) InsertOneStmt() (string, []any) {
	return "INSERT INTO c (id) VALUES (?);", v.Values()
}
func (v C) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM c WHERE id = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}

func (D) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"id"},
			Definition: "PRIMARY KEY (id)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "id", Definition: "id VARCHAR(255)"},
		},
	}
}
func (D) TableName() string {
	return "d"
}
func (D) HasPK() {}
func (v D) PK() (string, int, any) {
	return "id", 0, (driver.Valuer)(v.ID)
}
func (D) Columns() []string {
	return []string{"id"}
}
func (v D) Values() []any {
	return []any{(driver.Valuer)(v.ID)}
}
func (v *D) Addrs() []any {
	return []any{(sql.Scanner)(&v.ID)}
}
func (D) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v D) InsertOneStmt() (string, []any) {
	return "INSERT INTO d (id) VALUES (?);", v.Values()
}
func (v D) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM d WHERE id = ? LIMIT 1;", []any{(driver.Valuer)(v.ID)}
}
func (v D) GetID() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("id", v.ID, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
