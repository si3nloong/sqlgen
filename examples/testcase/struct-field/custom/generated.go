package custom

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Address) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`line_1`", Definition: "`line_1` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`line_2`", Definition: "`line_2` VARCHAR(255) DEFAULT ''"},
			{Name: "`city`", Definition: "`city` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`post_code`", Definition: "`post_code` INTEGER UNSIGNED NOT NULL DEFAULT 0"},
			{Name: "`state_code`", Definition: "`state_code` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`geo_point`", Definition: "`geo_point` VARCHAR(255) NOT NULL"},
			{Name: "`country_code`", Definition: "`country_code` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (Address) TableName() string {
	return "`address`"
}
func (Address) ColumnNames() []string {
	return []string{"`line_1`", "`line_2`", "`city`", "`post_code`", "`state_code`", "`geo_point`", "`country_code`"}
}
func (Address) SQLColumns() []string {
	return []string{"`line_1`", "`line_2`", "`city`", "`post_code`", "`state_code`", "ST_AsBinary(`geo_point`, 4326)", "`country_code`"}
}
func (v Address) Values() []any {
	return []any{string(v.Line1), (driver.Valuer)(v.Line2), string(v.City), int64(v.PostCode), string(v.StateCode), ewkb.Value(v.GeoPoint, 4326), string(v.CountryCode)}
}
func (v *Address) Addrs() []any {
	return []any{types.String(&v.Line1), (sql.Scanner)(&v.Line2), types.String(&v.City), types.Integer(&v.PostCode), types.String(&v.StateCode), ewkb.Scanner(&v.GeoPoint), types.String(&v.CountryCode)}
}
func (Address) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?)"
}
func (v Address) InsertOneStmt() (string, []any) {
	return "INSERT INTO `address` (`line_1`,`line_2`,`city`,`post_code`,`state_code`,`geo_point`,`country_code`) VALUES (?,?,?,?,?,ST_GeomFromEWKB(?),?);", v.Values()
}
func (v Address) GetLine1() sequel.ColumnValuer[string] {
	return sequel.Column("`line_1`", v.Line1, func(val string) driver.Value { return string(val) })
}
func (v Address) GetLine2() sequel.ColumnValuer[sql.NullString] {
	return sequel.Column("`line_2`", v.Line2, func(val sql.NullString) driver.Value { return (driver.Valuer)(val) })
}
func (v Address) GetCity() sequel.ColumnValuer[string] {
	return sequel.Column("`city`", v.City, func(val string) driver.Value { return string(val) })
}
func (v Address) GetPostCode() sequel.ColumnValuer[uint] {
	return sequel.Column("`post_code`", v.PostCode, func(val uint) driver.Value { return int64(val) })
}
func (v Address) GetStateCode() sequel.ColumnValuer[StateCode] {
	return sequel.Column("`state_code`", v.StateCode, func(val StateCode) driver.Value { return string(val) })
}
func (v Address) GetGeoPoint() sequel.SQLColumnValuer[orb.Point] {
	return sequel.SQLColumn("`geo_point`", v.GeoPoint, func(placeholder string) string { return "ST_GeomFromEWKB(" + placeholder + ")" }, func(val orb.Point) driver.Value { return ewkb.Value(val, 4326) })
}
func (v Address) GetCountryCode() sequel.ColumnValuer[CountryCode] {
	return sequel.Column("`country_code`", v.CountryCode, func(val CountryCode) driver.Value { return string(val) })
}

func (Customer) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT NOT NULL DEFAULT 0"},
			{Name: "`howOld`", Definition: "`howOld` TINYINT UNSIGNED NOT NULL DEFAULT 0"},
			{Name: "`name`", Definition: "`name` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`address`", Definition: "`address` JSON NOT NULL"},
			{Name: "`nicknames`", Definition: "`nicknames` JSON NOT NULL"},
			{Name: "`status`", Definition: "`status` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`join_at`", Definition: "`join_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"},
		},
	}
}
func (Customer) TableName() string {
	return "`customer`"
}
func (Customer) ColumnNames() []string {
	return []string{"`id`", "`howOld`", "`name`", "`address`", "`nicknames`", "`status`", "`join_at`"}
}
func (v Customer) Values() []any {
	return []any{int64(v.ID), int64(v.Age), (driver.Valuer)(v.Name), (driver.Valuer)(v.Address), encoding.MarshalStringList(v.Nicknames), string(v.Status), time.Time(v.JoinAt)}
}
func (v *Customer) Addrs() []any {
	return []any{types.Integer(&v.ID), types.Integer(&v.Age), (sql.Scanner)(&v.Name), types.JSONUnmarshaler(&v.Address), types.StringList(&v.Nicknames), types.String(&v.Status), (*time.Time)(&v.JoinAt)}
}
func (Customer) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?)"
}
func (v Customer) InsertOneStmt() (string, []any) {
	return "INSERT INTO `customer` (`id`,`howOld`,`name`,`address`,`nicknames`,`status`,`join_at`) VALUES (?,?,?,?,?,?,?);", v.Values()
}
func (v Customer) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("`id`", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v Customer) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("`howOld`", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v Customer) GetName() sequel.ColumnValuer[longText] {
	return sequel.Column("`name`", v.Name, func(val longText) driver.Value { return (driver.Valuer)(val) })
}
func (v Customer) GetAddress() sequel.ColumnValuer[Addresses] {
	return sequel.Column("`address`", v.Address, func(val Addresses) driver.Value { return (driver.Valuer)(val) })
}
func (v Customer) GetNicknames() sequel.ColumnValuer[[]longText] {
	return sequel.Column("`nicknames`", v.Nicknames, func(val []longText) driver.Value { return encoding.MarshalStringList(val) })
}
func (v Customer) GetStatus() sequel.ColumnValuer[string] {
	return sequel.Column("`status`", v.Status, func(val string) driver.Value { return string(val) })
}
func (v Customer) GetJoinAt() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`join_at`", v.JoinAt, func(val time.Time) driver.Value { return time.Time(val) })
}
