package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (AliasStruct) TableName() string {
	return "alias_struct"
}
func (AliasStruct) HasPK() {}
func (v AliasStruct) PK() (string, int, any) {
	return "Id", 1, v.ID
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"}
}
func (v AliasStruct) Values() []any {
	return []any{v.B, v.pk.ID, (string)(v.Header), string(v.Raw), (string)(v.Text), (driver.Valuer)(v.NullStr), (time.Time)(v.model.Created), (time.Time)(v.model.Updated)}
}
func (v *AliasStruct) Addrs() []any {
	return []any{&v.B, &v.pk.ID, encoding.StringScanner[aliasStr](&v.Header), encoding.StringScanner[sql.RawBytes](&v.Raw), encoding.StringScanner[customStr](&v.Text), (sql.Scanner)(&v.NullStr), encoding.TimeScanner(&v.model.Created), encoding.TimeScanner(&v.model.Updated)}
}
func (AliasStruct) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?)"
}
func (v AliasStruct) InsertOneStmt() (string, []any) {
	return "INSERT INTO alias_struct (b,Id,header,raw,text,null_str,created,updated) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v AliasStruct) FindOneByPKStmt() (string, []any) {
	return "SELECT b,Id,header,raw,text,null_str,created,updated FROM alias_struct WHERE Id = ? LIMIT 1;", []any{v.pk.ID}
}
func (v AliasStruct) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE alias_struct SET b = ?,header = ?,raw = ?,text = ?,null_str = ?,created = ?,updated = ? WHERE Id = ?;", []any{v.B, (string)(v.Header), string(v.Raw), (string)(v.Text), (driver.Valuer)(v.NullStr), (time.Time)(v.model.Created), (time.Time)(v.model.Updated), v.pk.ID}
}
func (v AliasStruct) GetB() driver.Value {
	return v.B
}
func (v AliasStruct) GetID() driver.Value {
	return v.ID
}
func (v AliasStruct) GetHeader() driver.Value {
	return (string)(v.Header)
}
func (v AliasStruct) GetRaw() driver.Value {
	return string(v.Raw)
}
func (v AliasStruct) GetText() driver.Value {
	return (string)(v.Text)
}
func (v AliasStruct) GetNullStr() driver.Value {
	return (driver.Valuer)(v.NullStr)
}
func (v AliasStruct) GetCreated() driver.Value {
	return (time.Time)(v.Created)
}
func (v AliasStruct) GetUpdated() driver.Value {
	return (time.Time)(v.Updated)
}

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"name"}
}
func (v B) Values() []any {
	return []any{v.Name}
}
func (v *B) Addrs() []any {
	return []any{&v.Name}
}
func (B) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v B) InsertOneStmt() (string, []any) {
	return "INSERT INTO b (name) VALUES (?);", v.Values()
}
func (v B) GetName() driver.Value {
	return v.Name
}

func (C) TableName() string {
	return "c"
}
func (C) Columns() []string {
	return []string{"id"}
}
func (v C) Values() []any {
	return []any{v.ID}
}
func (v *C) Addrs() []any {
	return []any{&v.ID}
}
func (C) InsertPlaceholders(row int) string {
	return "(?)"
}
func (v C) InsertOneStmt() (string, []any) {
	return "INSERT INTO c (id) VALUES (?);", v.Values()
}
func (v C) GetID() driver.Value {
	return v.ID
}
