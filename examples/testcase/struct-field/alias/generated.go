package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v AliasStruct) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (b FLOAT NOT NULL,Id BIGINT NOT NULL,header VARCHAR(255) NOT NULL,raw BLOB,text VARCHAR(255) NOT NULL,null_str VARCHAR(255) NOT NULL,created DATETIME NOT NULL,updated DATETIME NOT NULL,PRIMARY KEY (Id));"
}
func (AliasStruct) AlterTableStmt() string {
	return "ALTER TABLE alias_struct MODIFY b FLOAT NOT NULL,MODIFY Id BIGINT NOT NULL AFTER b,MODIFY header VARCHAR(255) NOT NULL AFTER Id,MODIFY raw BLOB AFTER header,MODIFY text VARCHAR(255) NOT NULL AFTER raw,MODIFY null_str VARCHAR(255) NOT NULL AFTER text,MODIFY created DATETIME NOT NULL AFTER null_str,MODIFY updated DATETIME NOT NULL AFTER created;"
}
func (AliasStruct) TableName() string {
	return "alias_struct"
}
func (v AliasStruct) InsertOneStmt() string {
	return "INSERT INTO alias_struct (b,Id,header,raw,text,null_str,created,updated) VALUES (?,?,?,?,?,?,?,?);"
}
func (AliasStruct) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?)"
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"}
}
func (v AliasStruct) PK() (columnName string, pos int, value driver.Value) {
	return "Id", 1, int64(v.pk.ID)
}
func (v AliasStruct) FindByPKStmt() string {
	return "SELECT b,Id,header,raw,text,null_str,created,updated FROM alias_struct WHERE Id = ? LIMIT 1;"
}
func (v AliasStruct) UpdateByPKStmt() string {
	return "UPDATE alias_struct SET b = ?,header = ?,raw = ?,text = ?,null_str = ?,created = ?,updated = ? WHERE Id = ? LIMIT 1;"
}
func (v AliasStruct) Values() []any {
	return []any{float64(v.B), int64(v.pk.ID), string(v.Header), string(v.Raw), string(v.Text), (driver.Valuer)(v.NullStr), time.Time(v.model.Created), time.Time(v.model.Updated)}
}
func (v *AliasStruct) Addrs() []any {
	return []any{types.Float(&v.B), types.Integer(&v.pk.ID), types.String(&v.Header), types.String(&v.Raw), types.String(&v.Text), (sql.Scanner)(&v.NullStr), (*time.Time)(&v.model.Created), (*time.Time)(&v.model.Updated)}
}
func (v AliasStruct) GetB() sequel.ColumnValuer[float64] {
	return sequel.Column[float64]("b", v.B, func(vi float64) driver.Value { return float64(vi) })
}
func (v AliasStruct) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("Id", v.pk.ID, func(vi int64) driver.Value { return int64(vi) })
}
func (v AliasStruct) GetHeader() sequel.ColumnValuer[customStr] {
	return sequel.Column[customStr]("header", v.Header, func(vi customStr) driver.Value { return string(vi) })
}
func (v AliasStruct) GetRaw() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column[sql.RawBytes]("raw", v.Raw, func(vi sql.RawBytes) driver.Value { return string(vi) })
}
func (v AliasStruct) GetText() sequel.ColumnValuer[customStr] {
	return sequel.Column[customStr]("text", v.Text, func(vi customStr) driver.Value { return string(vi) })
}
func (v AliasStruct) GetNullStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column[sql.NullString]("null_str", v.NullStr, func(vi sql.NullString) driver.Value { return (driver.Valuer)(vi) })
}
func (v AliasStruct) GetCreated() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("created", v.model.Created, func(vi time.Time) driver.Value { return time.Time(vi) })
}
func (v AliasStruct) GetUpdated() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("updated", v.model.Updated, func(vi time.Time) driver.Value { return time.Time(vi) })
}

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (name VARCHAR(255) NOT NULL);"
}
func (B) AlterTableStmt() string {
	return "ALTER TABLE b MODIFY name VARCHAR(255) NOT NULL;"
}
func (B) TableName() string {
	return "b"
}
func (v B) InsertOneStmt() string {
	return "INSERT INTO b (name) VALUES (?);"
}
func (B) InsertVarQuery() string {
	return "(?)"
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
func (v B) GetName() sequel.ColumnValuer[string] {
	return sequel.Column[string]("name", v.Name, func(vi string) driver.Value { return string(vi) })
}

func (v C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (id BIGINT NOT NULL);"
}
func (C) AlterTableStmt() string {
	return "ALTER TABLE c MODIFY id BIGINT NOT NULL;"
}
func (C) TableName() string {
	return "c"
}
func (v C) InsertOneStmt() string {
	return "INSERT INTO c (id) VALUES (?);"
}
func (C) InsertVarQuery() string {
	return "(?)"
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
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("id", v.ID, func(vi int64) driver.Value { return int64(vi) })
}
