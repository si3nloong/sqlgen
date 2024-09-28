package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (AliasStruct) TableName() string {
	return "alias_struct"
}
func (AliasStruct) HasPK() {}
func (v AliasStruct) PK() (string, int, any) {
	return "Id", 1, (int64)(v.pk.ID)
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"}
}
func (v AliasStruct) Values() []any {
	return []any{(float64)(v.B), (int64)(v.pk.ID), (string)(v.Header), string(v.Raw), (string)(v.Text), (driver.Valuer)(v.NullStr), (time.Time)(v.model.Created), (time.Time)(v.model.Updated)}
}
func (v *AliasStruct) Addrs() []any {
	return []any{types.Float64(&v.B), types.Integer(&v.pk.ID), types.String(&v.Header), types.String(&v.Raw), types.String(&v.Text), (sql.Scanner)(&v.NullStr), (*time.Time)(&v.model.Created), (*time.Time)(&v.model.Updated)}
}
func (AliasStruct) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?)"
}
func (v AliasStruct) InsertOneStmt() (string, []any) {
	return "INSERT INTO alias_struct (b,Id,header,raw,text,null_str,created,updated) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v AliasStruct) FindOneByPKStmt() (string, []any) {
	return "SELECT b,Id,header,raw,text,null_str,created,updated FROM alias_struct WHERE Id = ? LIMIT 1;", []any{(int64)(v.pk.ID)}
}
func (v AliasStruct) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE alias_struct SET b = ?,header = ?,raw = ?,text = ?,null_str = ?,created = ?,updated = ? WHERE Id = ?;", []any{(float64)(v.B), (string)(v.Header), string(v.Raw), (string)(v.Text), (driver.Valuer)(v.NullStr), (time.Time)(v.model.Created), (time.Time)(v.model.Updated), (int64)(v.pk.ID)}
}
func (v AliasStruct) GetB() sequel.ColumnValuer[float64] {
	return sequel.Column("b", v.B, func(val float64) driver.Value { return (float64)(val) })
}
func (v AliasStruct) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("Id", v.pk.ID, func(val int64) driver.Value { return (int64)(val) })
}
func (v AliasStruct) GetHeader() sequel.ColumnValuer[aliasStr] {
	return sequel.Column("header", v.Header, func(val aliasStr) driver.Value { return (string)(val) })
}
func (v AliasStruct) GetRaw() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw", v.Raw, func(val sql.RawBytes) driver.Value { return string(val) })
}
func (v AliasStruct) GetText() sequel.ColumnValuer[customStr] {
	return sequel.Column("text", v.Text, func(val customStr) driver.Value { return (string)(val) })
}
func (v AliasStruct) GetNullStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("null_str", v.NullStr, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
func (v AliasStruct) GetCreated() sequel.ColumnValuer[DT] {
	return sequel.Column("created", v.model.Created, func(val DT) driver.Value { return (time.Time)(val) })
}
func (v AliasStruct) GetUpdated() sequel.ColumnValuer[DT] {
	return sequel.Column("updated", v.model.Updated, func(val DT) driver.Value { return (time.Time)(val) })
}

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"name"}
}
func (v B) Values() []any {
	return []any{(string)(v.Name)}
}
func (v *B) Addrs() []any {
	return []any{types.String(&v.Name)}
}
func (B) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (name) VALUES (?);", v.Values()
}
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value { return (string)(val) })
}

func (C) TableName() string {
	return "c"
}
func (C) Columns() []string {
	return []string{"id"}
}
func (v C) Values() []any {
	return []any{(int64)(v.ID)}
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
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return (int64)(val) })
}
