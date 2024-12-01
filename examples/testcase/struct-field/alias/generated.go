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
	return "Id", 1, v.pk.ID
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"} // 8
}
func (v AliasStruct) Values() []any {
	return []any{
		v.B,                          // 0 - b
		v.pk.ID,                      // 1 - Id
		(string)(v.Header),           // 2 - header
		string(v.Raw),                // 3 - raw
		(string)(v.Text),             // 4 - text
		v.NullStr,                    // 5 - null_str
		(time.Time)(v.model.Created), // 6 - created
		(time.Time)(v.model.Updated), // 7 - updated
	}
}
func (v *AliasStruct) Addrs() []any {
	return []any{
		&v.B,     // 0 - b
		&v.pk.ID, // 1 - Id
		encoding.StringScanner[aliasStr](&v.Header),  // 2 - header
		encoding.StringScanner[sql.RawBytes](&v.Raw), // 3 - raw
		encoding.StringScanner[customStr](&v.Text),   // 4 - text
		&v.NullStr,                             // 5 - null_str
		encoding.TimeScanner(&v.model.Created), // 6 - created
		encoding.TimeScanner(&v.model.Updated), // 7 - updated
	}
}
func (AliasStruct) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?)" // 8
}
func (v AliasStruct) InsertOneStmt() (string, []any) {
	return "INSERT INTO alias_struct (b,Id,header,raw,text,null_str,created,updated) VALUES (?,?,?,?,?,?,?,?);", v.Values()
}
func (v AliasStruct) FindOneByPKStmt() (string, []any) {
	return "SELECT b,Id,header,raw,text,null_str,created,updated FROM alias_struct WHERE Id = ? LIMIT 1;", []any{v.pk.ID}
}
func (v AliasStruct) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE alias_struct SET b = ?,header = ?,raw = ?,text = ?,null_str = ?,created = ?,updated = ? WHERE Id = ?;", []any{v.B, (string)(v.Header), string(v.Raw), (string)(v.Text), v.NullStr, (time.Time)(v.model.Created), (time.Time)(v.model.Updated), v.pk.ID}
}
func (v AliasStruct) GetB() driver.Value {
	return v.B
}
func (v AliasStruct) GetID() driver.Value {
	return v.pk.ID
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
	return v.NullStr
}
func (v AliasStruct) GetCreated() driver.Value {
	return (time.Time)(v.model.Created)
}
func (v AliasStruct) GetUpdated() driver.Value {
	return (time.Time)(v.model.Updated)
}

func (B) TableName() string {
	return "b"
}
func (B) Columns() []string {
	return []string{"name"} // 1
}
func (v B) Values() []any {
	return []any{
		v.Name, // 0 - name
	}
}
func (v *B) Addrs() []any {
	return []any{
		&v.Name, // 0 - name
	}
}
func (B) InsertPlaceholders(row int) string {
	return "(?)" // 1
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
func (v C) GetID() driver.Value {
	return v.ID
}
