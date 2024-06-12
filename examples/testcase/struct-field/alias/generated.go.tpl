package alias

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v AliasStruct) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `alias_struct` (`b` FLOAT NOT NULL,`Id` BIGINT NOT NULL,`header` VARCHAR(255) NOT NULL,`raw` BLOB NOT NULL,`text` VARCHAR(255) NOT NULL,`null_str` VARCHAR(255) NOT NULL,`created` DATETIME NOT NULL,`updated` DATETIME NOT NULL,PRIMARY KEY (`Id`));"
}
func (AliasStruct) TableName() string {
	return "alias_struct"
}
func (AliasStruct) InsertOneStmt() string {
	return "INSERT INTO alias_struct (b,Id,header,raw,text,null_str,created,updated) VALUES (?,?,?,?,?,?,?,?);"
}
func (AliasStruct) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?,?)"
}
func (AliasStruct) Columns() []string {
	return []string{"b", "Id", "header", "raw", "text", "null_str", "created", "updated"}
}
func (AliasStruct) HasPK() {}
func (v AliasStruct) PK() (string, int, any) {
	return "Id", 1, int64(v.pk.ID)
}
func (AliasStruct) FindByPKStmt() string {
	return "SELECT b,Id,header,raw,text,null_str,created,updated FROM alias_struct WHERE Id = ? LIMIT 1;"
}
func (AliasStruct) UpdateByPKStmt() string {
	return "UPDATE alias_struct SET b = ?,header = ?,raw = ?,text = ?,null_str = ?,created = ?,updated = ? WHERE Id = ? LIMIT 1;"
}
func (v AliasStruct) Values() []any {
	return []any{float64(v.B), int64(v.pk.ID), string(v.Header), string(v.Raw), string(v.Text), (driver.Valuer)(v.NullStr), time.Time(v.model.Created), time.Time(v.model.Updated)}
}
func (v *AliasStruct) Addrs() []any {
	return []any{types.Float(&v.B), types.Integer(&v.pk.ID), types.String(&v.Header), types.String(&v.Raw), types.String(&v.Text), (sql.Scanner)(&v.NullStr), (*time.Time)(&v.model.Created), (*time.Time)(&v.model.Updated)}
}
func (v AliasStruct) GetB() sequel.ColumnValuer[float64] {
	return sequel.Column("b", v.B, func(val float64) driver.Value { return float64(val) })
}
func (v AliasStruct) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("Id", v.pk.ID, func(val int64) driver.Value { return int64(val) })
}
func (v AliasStruct) GetHeader() sequel.ColumnValuer[customStr] {
	return sequel.Column("header", v.Header, func(val customStr) driver.Value { return string(val) })
}
func (v AliasStruct) GetRaw() sequel.ColumnValuer[sql.RawBytes] {
	return sequel.Column("raw", v.Raw, func(val sql.RawBytes) driver.Value { return string(val) })
}
func (v AliasStruct) GetText() sequel.ColumnValuer[customStr] {
	return sequel.Column("text", v.Text, func(val customStr) driver.Value { return string(val) })
}
func (v AliasStruct) GetNullStr() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("null_str", v.NullStr, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
func (v AliasStruct) GetCreated() sequel.ColumnValuer[time.Time] {
	return sequel.Column("created", v.model.Created, func(val time.Time) driver.Value { return time.Time(val) })
}
func (v AliasStruct) GetUpdated() sequel.ColumnValuer[time.Time] {
	return sequel.Column("updated", v.model.Updated, func(val time.Time) driver.Value { return time.Time(val) })
}

func (v B) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `b` (`name` VARCHAR(255) NOT NULL);"
}
func (B) TableName() string {
	return "b"
}
func (B) InsertOneStmt() string {
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
	return sequel.Column("name", v.Name, func(val string) driver.Value { return string(val) })
}

func (v C) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `c` (`id` BIGINT NOT NULL);"
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
func (v C) Values() []any {
	return []any{int64(v.ID)}
}
func (v *C) Addrs() []any {
	return []any{types.Integer(&v.ID)}
}
func (v C) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}
