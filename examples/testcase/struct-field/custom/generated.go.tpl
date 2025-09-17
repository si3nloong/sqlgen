package custom

import (
	"database/sql"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/sqltype"
)

func (Address) TableName() string {
	return "address"
}
func (Address) SQLColumns() []string {
	return []string{"line_1", "line_2", "city", "post_code", "state_code", "ST_AsBinary(geo_point,4326)", "country_code"} // 7
}
func (Address) Columns() []string {
	return []string{"line_1", "line_2", "city", "post_code", "state_code", "geo_point", "country_code"} // 7
}
func (v Address) Values() []any {
	return []any{
		v.Line1,                      // 0 - line_1
		v.Line2,                      // 1 - line_2
		v.City,                       // 2 - city
		(int64)(v.PostCode),          // 3 - post_code
		(string)(v.StateCode),        // 4 - state_code
		ewkb.Value(v.GeoPoint, 4326), // 5 - geo_point
		(string)(v.CountryCode),      // 6 - country_code
	}
}
func (v *Address) Addrs() []any {
	return []any{
		&v.Line1,                                // 0 - line_1
		&v.Line2,                                // 1 - line_2
		&v.City,                                 // 2 - city
		encoding.UintScanner[uint](&v.PostCode), // 3 - post_code
		encoding.StringScanner[StateCode](&v.StateCode),     // 4 - state_code
		ewkb.Scanner(&v.GeoPoint),                           // 5 - geo_point
		encoding.StringScanner[CountryCode](&v.CountryCode), // 6 - country_code
	}
}
func (Address) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?)" // 7
}
func (v Address) InsertOneStmt() (string, []any) {
	return "INSERT INTO `address` (`line_1`,`line_2`,`city`,`post_code`,`state_code`,`geo_point`,`country_code`) VALUES (?,?,?,?,?,ST_GeomFromEWKB(?),?);", v.Values()
}
func (v Address) Line1Value() any {
	return v.Line1
}
func (v Address) Line2Value() any {
	return v.Line2
}
func (v Address) CityValue() any {
	return v.City
}
func (v Address) PostCodeValue() any {
	return (int64)(v.PostCode)
}
func (v Address) StateCodeValue() any {
	return (string)(v.StateCode)
}
func (v Address) GeoPointValue() any {
	return ewkb.Value(v.GeoPoint, 4326)
}
func (v Address) CountryCodeValue() any {
	return (string)(v.CountryCode)
}
func (v Address) ColumnLine1() sequel.ColumnClause[string] {
	return sequel.BasicColumn("line_1", v.Line1)
}
func (v Address) ColumnLine2() sequel.ColumnConvertClause[sql.NullString] {
	return sequel.Column("line_2", v.Line2, func(val sql.NullString) any {
		return val
	})
}
func (v Address) ColumnCity() sequel.ColumnClause[string] {
	return sequel.BasicColumn("city", v.City)
}
func (v Address) ColumnPostCode() sequel.ColumnConvertClause[uint] {
	return sequel.Column("post_code", v.PostCode, func(val uint) any {
		return (int64)(val)
	})
}
func (v Address) ColumnStateCode() sequel.ColumnConvertClause[StateCode] {
	return sequel.Column("state_code", v.StateCode, func(val StateCode) any {
		return (string)(val)
	})
}
func (v Address) ColumnGeoPoint() sequel.ColumnConvertClause[orb.Point] {
	return sequel.Column("geo_point", v.GeoPoint, func(val orb.Point) any {
		return ewkb.Value(val, 4326)
	})
}
func (v Address) ColumnCountryCode() sequel.ColumnConvertClause[CountryCode] {
	return sequel.Column("country_code", v.CountryCode, func(val CountryCode) any {
		return (string)(val)
	})
}

func (Customer) TableName() string {
	return "customer"
}
func (Customer) Columns() []string {
	return []string{"id", "howOld", "name", "address", "nicknames", "status", "join_at"} // 7
}
func (v Customer) Values() []any {
	return []any{
		v.ID,           // 0 - id
		(int64)(v.Age), // 1 - howOld
		v.Name,         // 2 - name
		v.Address,      // 3 - address
		(sqltype.StringSlice[longText])(v.Nicknames), // 4 - nicknames
		v.Status, // 5 - status
		v.JoinAt, // 6 - join_at
	}
}
func (v *Customer) Addrs() []any {
	return []any{
		&v.ID,                                // 0 - id
		encoding.Uint8Scanner[uint8](&v.Age), // 1 - howOld
		&v.Name,                              // 2 - name
		encoding.JSONScanner(&v.Address),     // 3 - address
		(*sqltype.StringSlice[longText])(&v.Nicknames), // 4 - nicknames
		&v.Status, // 5 - status
		&v.JoinAt, // 6 - join_at
	}
}
func (Customer) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?)" // 7
}
func (v Customer) InsertOneStmt() (string, []any) {
	return "INSERT INTO `customer` (`id`,`howOld`,`name`,`address`,`nicknames`,`status`,`join_at`) VALUES (?,?,?,?,?,?,?);", v.Values()
}
func (v Customer) IDValue() any {
	return v.ID
}
func (v Customer) AgeValue() any {
	return (int64)(v.Age)
}
func (v Customer) NameValue() any {
	return v.Name
}
func (v Customer) AddressValue() any {
	return v.Address
}
func (v Customer) NicknamesValue() any {
	return (sqltype.StringSlice[longText])(v.Nicknames)
}
func (v Customer) StatusValue() any {
	return v.Status
}
func (v Customer) JoinAtValue() any {
	return v.JoinAt
}
func (v Customer) ColumnID() sequel.ColumnClause[int64] {
	return sequel.BasicColumn("id", v.ID)
}
func (v Customer) ColumnAge() sequel.ColumnConvertClause[uint8] {
	return sequel.Column("howOld", v.Age, func(val uint8) any {
		return (int64)(val)
	})
}
func (v Customer) ColumnName() sequel.ColumnConvertClause[longText] {
	return sequel.Column("name", v.Name, func(val longText) any {
		return val
	})
}
func (v Customer) ColumnAddress() sequel.ColumnConvertClause[Addresses] {
	return sequel.Column("address", v.Address, func(val Addresses) any {
		return val
	})
}
func (v Customer) ColumnNicknames() sequel.ColumnConvertClause[[]longText] {
	return sequel.Column("nicknames", v.Nicknames, func(val []longText) any {
		return (sqltype.StringSlice[longText])(val)
	})
}
func (v Customer) ColumnStatus() sequel.ColumnClause[string] {
	return sequel.BasicColumn("status", v.Status)
}
func (v Customer) ColumnJoinAt() sequel.ColumnClause[time.Time] {
	return sequel.BasicColumn("join_at", v.JoinAt)
}
