package schema

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (A) TableName() string {
	return "Apple"
}
func (A) Columns() []string {
	return []string{"id", "text", "created_at"} // 3
}
func (v A) Values() []any {
	return []any{
		v.ID,             // 0 - id
		(string)(v.Text), // 1 - text
		v.CreatedAt,      // 2 - created_at
	}
}
func (v *A) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
		encoding.StringScanner[LongText](&v.Text), // 1 - text
		&v.CreatedAt, // 2 - created_at
	}
}
func (A) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v A) InsertOneStmt() (string, []any) {
	return "INSERT INTO Apple (id,text,created_at) VALUES (?,?,?);", v.Values()
}
func (v A) GetID() driver.Value {
	return v.ID
}
func (v A) GetText() driver.Value {
	return (string)(v.Text)
}
func (v A) GetCreatedAt() driver.Value {
	return v.CreatedAt
}

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"id", "created_at"} // 2
}
func (v B) Values() []any {
	return []any{
		v.ID,        // 0 - id
		v.CreatedAt, // 1 - created_at
	}
}
func (v *B) Addrs() []any {
	return []any{
		&v.ID,        // 0 - id
		&v.CreatedAt, // 1 - created_at
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (id,created_at) VALUES (?,?);", v.Values()
}
func (v B) GetID() driver.Value {
	return v.ID
}
func (v B) GetCreatedAt() driver.Value {
	return v.CreatedAt
}

func (C) TableName() string {
	return "c"
}
func (C) HasPK() {}
func (v C) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (C) Columns() []string {
	return []string{"id"} // 1
}
func (v C) Values() []any {
	return []any{
		v.ID, // 0 - id
	}
}
func (v *C) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
	}
}
func (C) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v C) InsertOneStmt() (string, []any) {
	return "INSERT INTO c (id) VALUES (?);", v.Values()
}
func (v C) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM c WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v C) GetID() driver.Value {
	return v.ID
}

func (D) TableName() string {
	return "d"
}
func (D) HasPK() {}
func (v D) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (D) Columns() []string {
	return []string{"id"} // 1
}
func (v D) Values() []any {
	return []any{
		v.ID, // 0 - id
	}
}
func (v *D) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
	}
}
func (D) InsertPlaceholders(row int) string {
	return "(?)" // 1
}
func (v D) InsertOneStmt() (string, []any) {
	return "INSERT INTO d (id) VALUES (?);", v.Values()
}
func (v D) FindOneByPKStmt() (string, []any) {
	return "SELECT id FROM d WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v D) GetID() driver.Value {
	return v.ID
}
