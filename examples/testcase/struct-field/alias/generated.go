package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
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
func (v AliasStruct) BValue() driver.Value {
	return v.B
}
func (v AliasStruct) IDValue() driver.Value {
	return v.pk.ID
}
func (v AliasStruct) HeaderValue() driver.Value {
	return (string)(v.Header)
}
func (v AliasStruct) RawValue() driver.Value {
	return string(v.Raw)
}
func (v AliasStruct) TextValue() driver.Value {
	return (string)(v.Text)
}
func (v AliasStruct) NullStrValue() driver.Value {
	return v.NullStr
}
func (v AliasStruct) CreatedValue() driver.Value {
	return (time.Time)(v.model.Created)
}
func (v AliasStruct) UpdatedValue() driver.Value {
	return (time.Time)(v.model.Updated)
}
func (v AliasStruct) GetB() sequel.ColumnValuer[float64] {
	return sequel.Column("b", v.B, func(val float64) driver.Value {
		return val
	})
}
func (v AliasStruct) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("Id", v.pk.ID, func(val int64) driver.Value {
		return val
	})
}
func (v AliasStruct) GetHeader() sequel.ColumnValuer[aliasStr] {
	return sequel.Column("header", v.Header, func(val aliasStr) driver.Value {
		return (string)(val)
	})
}
func (v AliasStruct) GetRaw() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw", v.Raw, func(val sql.RawBytes) driver.Value {
		return string(val)
	})
}
func (v AliasStruct) GetText() sequel.ColumnValuer[customStr] {
	return sequel.Column("text", v.Text, func(val customStr) driver.Value {
		return (string)(val)
	})
}
func (v AliasStruct) GetNullStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("null_str", v.NullStr, func(val sql.NullString) driver.Value {
		return val
	})
}
func (v AliasStruct) GetCreated() sequel.ColumnValuer[DT] {
	return sequel.Column("created", v.model.Created, func(val DT) driver.Value {
		return (time.Time)(val)
	})
}
func (v AliasStruct) GetUpdated() sequel.ColumnValuer[DT] {
	return sequel.Column("updated", v.model.Updated, func(val DT) driver.Value {
		return (time.Time)(val)
	})
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
func (v B) NameValue() driver.Value {
	return v.Name
}
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column("name", v.Name, func(val string) driver.Value {
		return val
	})
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
func (v C) IDValue() driver.Value {
	return v.ID
}
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value {
		return val
	})
}
