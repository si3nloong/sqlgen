package custom

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Address) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (line_1 VARCHAR(255) NOT NULL,line_2 VARCHAR(255) NOT NULL,city VARCHAR(255) NOT NULL,post_code INTEGER UNSIGNED NOT NULL,state_code VARCHAR(255) NOT NULL,country_code VARCHAR(255) NOT NULL);"
}
func (Address) AlterTableStmt() string {
	return "ALTER TABLE address MODIFY line_1 VARCHAR(255) NOT NULL,MODIFY line_2 VARCHAR(255) NOT NULL AFTER line_1,MODIFY city VARCHAR(255) NOT NULL AFTER line_2,MODIFY post_code INTEGER UNSIGNED NOT NULL AFTER city,MODIFY state_code VARCHAR(255) NOT NULL AFTER post_code,MODIFY country_code VARCHAR(255) NOT NULL AFTER state_code;"
}
func (Address) TableName() string {
	return "address"
}
func (v Address) InsertOneStmt() string {
	return "INSERT INTO address (line_1,line_2,city,post_code,state_code,country_code) VALUES (?,?,?,?,?,?);"
}
func (Address) InsertVarQuery() string {
	return "(?,?,?,?,?,?)"
}
func (Address) Columns() []string {
	return []string{"line_1", "line_2", "city", "post_code", "state_code", "country_code"}
}
func (v Address) Values() []any {
	return []any{string(v.Line1), (driver.Valuer)(v.Line2), string(v.City), int64(v.PostCode), string(v.StateCode), string(v.CountryCode)}
}
func (v *Address) Addrs() []any {
	return []any{types.String(&v.Line1), (sql.Scanner)(&v.Line2), types.String(&v.City), types.Integer(&v.PostCode), types.String(&v.StateCode), types.String(&v.CountryCode)}
}
func (v Address) GetLine1() sequel.ColumnValuer[string] {
	return sequel.Column[string]("line_1", v.Line1, func(vi string) driver.Value { return string(vi) })
}
func (v Address) GetLine2() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column[sql.NullString]("line_2", v.Line2, func(vi sql.NullString) driver.Value { return (driver.Valuer)(vi) })
}
func (v Address) GetCity() sequel.ColumnValuer[string] {
	return sequel.Column[string]("city", v.City, func(vi string) driver.Value { return string(vi) })
}
func (v Address) GetPostCode() sequel.ColumnValuer[uint] {
	return sequel.Column[uint]("post_code", v.PostCode, func(vi uint) driver.Value { return int64(vi) })
}
func (v Address) GetStateCode() sequel.ColumnValuer[StateCode] {
	return sequel.Column[StateCode]("state_code", v.StateCode, func(vi StateCode) driver.Value { return string(vi) })
}
func (v Address) GetCountryCode() sequel.ColumnValuer[CountryCode] {
	return sequel.Column[CountryCode]("country_code", v.CountryCode, func(vi CountryCode) driver.Value { return string(vi) })
}

func (v Customer) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (id BIGINT NOT NULL,howOld TINYINT UNSIGNED NOT NULL,name VARCHAR(255) NOT NULL,address JSON NOT NULL,nicknames JSON NOT NULL,status VARCHAR(255) NOT NULL,join_at DATETIME NOT NULL);"
}
func (Customer) AlterTableStmt() string {
	return "ALTER TABLE customer MODIFY id BIGINT NOT NULL,MODIFY howOld TINYINT UNSIGNED NOT NULL AFTER id,MODIFY name VARCHAR(255) NOT NULL AFTER howOld,MODIFY address JSON NOT NULL AFTER name,MODIFY nicknames JSON NOT NULL AFTER address,MODIFY status VARCHAR(255) NOT NULL AFTER nicknames,MODIFY join_at DATETIME NOT NULL AFTER status;"
}
func (Customer) TableName() string {
	return "customer"
}
func (v Customer) InsertOneStmt() string {
	return "INSERT INTO customer (id,howOld,name,address,nicknames,status,join_at) VALUES (?,?,?,?,?,?,?);"
}
func (Customer) InsertVarQuery() string {
	return "(?,?,?,?,?,?,?)"
}
func (Customer) Columns() []string {
	return []string{"id", "howOld", "name", "address", "nicknames", "status", "join_at"}
}
func (v Customer) Values() []any {
	return []any{int64(v.ID), int64(v.Age), (driver.Valuer)(v.Name), (driver.Valuer)(v.Address), encoding.MarshalStringList(v.Nicknames), string(v.Status), time.Time(v.JoinAt)}
}
func (v *Customer) Addrs() []any {
	return []any{types.Integer(&v.ID), types.Integer(&v.Age), (sql.Scanner)(&v.Name), &v.Address, types.StringList(&v.Nicknames), types.String(&v.Status), (*time.Time)(&v.JoinAt)}
}
func (v Customer) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("id", v.ID, func(vi int64) driver.Value { return int64(vi) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column[uint8]("howOld", v.Age, func(vi uint8) driver.Value { return int64(vi) })
}
func (v Customer) GetName() sequel.ColumnValuer[longText] {
	return sequel.Column[longText]("name", v.Name, func(vi longText) driver.Value { return (driver.Valuer)(vi) })
}
func (v Customer) GetAddress() sequel.ColumnValuer[Addresses] {
	return sequel.Column[Addresses]("address", v.Address, func(vi Addresses) driver.Value { return (driver.Valuer)(vi) })
}
func (v Customer) GetNicknames() sequel.ColumnValuer[[]longText] {
	return sequel.Column[[]longText]("nicknames", v.Nicknames, func(vi []longText) driver.Value { return encoding.MarshalStringList(vi) })
}
func (v Customer) GetStatus() sequel.ColumnValuer[string] {
	return sequel.Column[string]("status", v.Status, func(vi string) driver.Value { return string(vi) })
}
func (v Customer) GetJoinAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("join_at", v.JoinAt, func(vi time.Time) driver.Value { return time.Time(vi) })
}
