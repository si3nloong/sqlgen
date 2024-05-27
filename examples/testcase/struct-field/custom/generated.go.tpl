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
	return "CREATE TABLE IF NOT EXISTS `address` (`line_1` VARCHAR(255) NOT NULL,`line_2` VARCHAR(255) NOT NULL,`city` VARCHAR(255) NOT NULL,`post_code` INTEGER NOT NULL,`state_code` VARCHAR(255) NOT NULL,`country_code` VARCHAR(255) NOT NULL);"
}
func (Address) TableName() string {
	return "address"
}
func (Address) InsertOneStmt() string {
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
	return sequel.Column("line_1", v.Line1, func(val string) driver.Value { return string(val) })
}
func (v Address) GetLine2() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("line_2", v.Line2, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
func (v Address) GetCity() sequel.ColumnValuer[string] {
	return sequel.Column("city", v.City, func(val string) driver.Value { return string(val) })
}
func (v Address) GetPostCode() sequel.ColumnValuer[uint] {
	return sequel.Column("post_code", v.PostCode, func(val uint) driver.Value { return int64(val) })
}
func (v Address) GetStateCode() sequel.ColumnValuer[StateCode] {
	return sequel.Column("state_code", v.StateCode, func(val StateCode) driver.Value { return string(val) })
}
func (v Address) GetCountryCode() sequel.ColumnValuer[CountryCode] {
	return sequel.Column("country_code", v.CountryCode, func(val CountryCode) driver.Value { return string(val) })
}

func (v Customer) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `customer` (`id` BIGINT NOT NULL,`howOld` TINYINT UNSIGNED NOT NULL,`name` VARCHAR(255) NOT NULL,`address` JSON NOT NULL,`nicknames` JSON NOT NULL,`status` VARCHAR(255) NOT NULL,`join_at` DATETIME NOT NULL);"
}
func (Customer) TableName() string {
	return "customer"
}
func (Customer) InsertOneStmt() string {
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
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("howOld", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v Customer) GetName() sequel.ColumnValuer[longText] {
	return sequel.Column("name", v.Name, func(val longText) driver.Value { return (driver.Valuer)(val) })
}
func (v Customer) GetAddress() sequel.ColumnValuer[Addresses] {
	return sequel.Column("address", v.Address, func(val Addresses) driver.Value { return (driver.Valuer)(val) })
}
func (v Customer) GetNicknames() sequel.ColumnValuer[[]longText] {
	return sequel.Column("nicknames", v.Nicknames, func(val []longText) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Customer) GetStatus() sequel.ColumnValuer[string] {
	return sequel.Column("status", v.Status, func(val string) driver.Value { return string(val) })
}
func (v Customer) GetJoinAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column("join_at", v.JoinAt, func(val time.Time) driver.Value { return time.Time(val) })
}
