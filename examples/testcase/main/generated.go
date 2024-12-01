package main

import (
	"database/sql/driver"
	"reflect"

	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/sqltype"
)

func (Address) TableName() string {
	return "address"
}
func (Address) Columns() []string {
	return []string{"line_1", "line_2", "country_code"} // 3
}
func (v Address) Values() []any {
	return []any{
		v.Line1,       // 0 - line_1
		v.Line2,       // 1 - line_2
		v.CountryCode, // 2 - country_code
	}
}
func (v *Address) Addrs() []any {
	return []any{
		&v.Line1,       // 0 - line_1
		&v.Line2,       // 1 - line_2
		&v.CountryCode, // 2 - country_code
	}
}
func (Address) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Address) InsertOneStmt() (string, []any) {
	return "INSERT INTO address (line_1,line_2,country_code) VALUES (?,?,?);", v.Values()
}
func (v Address) GetLine1() driver.Value {
	return v.Line1
}
func (v Address) GetLine2() driver.Value {
	return v.Line2
}
func (v Address) GetCountryCode() driver.Value {
	return v.CountryCode
}

func (HouseUnit) TableName() string {
	return "house_unit"
}
func (HouseUnit) Columns() []string {
	return []string{"no", "build_time", "address", "kind", "type", "chan", "inner", "arr", "slice", "map"} // 10
}
func (v HouseUnit) Values() []any {
	return []any{
		(int64)(v.No),                            //  0 - no
		v.BuildTime,                              //  1 - build_time
		encoding.JSONValue(v.Address),            //  2 - address
		(int64)(v.Kind),                          //  3 - kind
		(int64)(v.Type),                          //  4 - type
		(int64)(v.Chan),                          //  5 - chan
		encoding.JSONValue(v.Inner),              //  6 - inner
		encoding.JSONValue(v.Arr),                //  7 - arr
		(sqltype.Float64Slice[float64])(v.Slice), //  8 - slice
		encoding.JSONValue(v.Map),                //  9 - map
	}
}
func (v *HouseUnit) Addrs() []any {
	return []any{
		encoding.UintScanner[uint](&v.No),             //  0 - no
		&v.BuildTime,                                  //  1 - build_time
		encoding.JSONScanner(&v.Address),              //  2 - address
		encoding.UintScanner[reflect.Kind](&v.Kind),   //  3 - kind
		encoding.Uint8Scanner[HouseUnitType](&v.Type), //  4 - type
		encoding.IntScanner[reflect.ChanDir](&v.Chan), //  5 - chan
		encoding.JSONScanner(&v.Inner),                //  6 - inner
		encoding.JSONScanner(&v.Arr),                  //  7 - arr
		(*sqltype.Float64Slice[float64])(&v.Slice),    //  8 - slice
		encoding.JSONScanner(&v.Map),                  //  9 - map
	}
}
func (HouseUnit) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?)" // 10
}
func (v HouseUnit) InsertOneStmt() (string, []any) {
	return "INSERT INTO house_unit (no,build_time,address,kind,type,chan,inner,arr,slice,map) VALUES (?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v HouseUnit) GetNo() driver.Value {
	return (int64)(v.No)
}
func (v HouseUnit) GetBuildTime() driver.Value {
	return v.BuildTime
}
func (v HouseUnit) GetAddress() driver.Value {
	return encoding.JSONValue(v.Address)
}
func (v HouseUnit) GetKind() driver.Value {
	return (int64)(v.Kind)
}
func (v HouseUnit) GetType() driver.Value {
	return (int64)(v.Type)
}
func (v HouseUnit) GetChan() driver.Value {
	return (int64)(v.Chan)
}
func (v HouseUnit) GetInner() driver.Value {
	return encoding.JSONValue(v.Inner)
}
func (v HouseUnit) GetArr() driver.Value {
	return encoding.JSONValue(v.Arr)
}
func (v HouseUnit) GetSlice() driver.Value {
	return (sqltype.Float64Slice[float64])(v.Slice)
}
func (v HouseUnit) GetMap() driver.Value {
	return encoding.JSONValue(v.Map)
}
