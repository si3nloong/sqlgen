package main

import (
	"database/sql/driver"
	"reflect"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
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
	return "INSERT INTO `address` (`line_1`,`line_2`,`country_code`) VALUES (?,?,?);", v.Values()
}
func (v Address) Line1Value() any {
	return v.Line1
}
func (v Address) Line2Value() any {
	return v.Line2
}
func (v Address) CountryCodeValue() any {
	return v.CountryCode
}
func (v Address) ColumnLine1() sequel.ColumnValuer[string] {
	return sequel.Column("line_1", v.Line1, func(val string) driver.Value {
		return val
	})
}
func (v Address) ColumnLine2() sequel.ColumnValuer[string] {
	return sequel.Column("line_2", v.Line2, func(val string) driver.Value {
		return val
	})
}
func (v Address) ColumnCountryCode() sequel.ColumnValuer[string] {
	return sequel.Column("country_code", v.CountryCode, func(val string) driver.Value {
		return val
	})
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
	return "INSERT INTO `house_unit` (`no`,`build_time`,`address`,`kind`,`type`,`chan`,`inner`,`arr`,`slice`,`map`) VALUES (?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v HouseUnit) NoValue() any {
	return (int64)(v.No)
}
func (v HouseUnit) BuildTimeValue() any {
	return v.BuildTime
}
func (v HouseUnit) AddressValue() any {
	return encoding.JSONValue(v.Address)
}
func (v HouseUnit) KindValue() any {
	return (int64)(v.Kind)
}
func (v HouseUnit) TypeValue() any {
	return (int64)(v.Type)
}
func (v HouseUnit) ChanValue() any {
	return (int64)(v.Chan)
}
func (v HouseUnit) InnerValue() any {
	return encoding.JSONValue(v.Inner)
}
func (v HouseUnit) ArrValue() any {
	return encoding.JSONValue(v.Arr)
}
func (v HouseUnit) SliceValue() any {
	return (sqltype.Float64Slice[float64])(v.Slice)
}
func (v HouseUnit) MapValue() any {
	return encoding.JSONValue(v.Map)
}
func (v HouseUnit) ColumnNo() sequel.ColumnValuer[uint] {
	return sequel.Column("no", v.No, func(val uint) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) ColumnBuildTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("build_time", v.BuildTime, func(val time.Time) driver.Value {
		return val
	})
}
func (v HouseUnit) ColumnAddress() sequel.ColumnValuer[Address] {
	return sequel.Column("address", v.Address, func(val Address) driver.Value {
		return encoding.JSONValue(val)
	})
}
func (v HouseUnit) ColumnKind() sequel.ColumnValuer[reflect.Kind] {
	return sequel.Column("kind", v.Kind, func(val reflect.Kind) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) ColumnType() sequel.ColumnValuer[HouseUnitType] {
	return sequel.Column("type", v.Type, func(val HouseUnitType) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) ColumnChan() sequel.ColumnValuer[reflect.ChanDir] {
	return sequel.Column("chan", v.Chan, func(val reflect.ChanDir) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) ColumnInner() sequel.ColumnValuer[struct{ Flag bool }] {
	return sequel.Column("inner", v.Inner, func(val struct{ Flag bool }) driver.Value {
		return encoding.JSONValue(val)
	})
}
func (v HouseUnit) ColumnArr() sequel.ColumnValuer[[2]string] {
	return sequel.Column("arr", v.Arr, func(val [2]string) driver.Value {
		return encoding.JSONValue(val)
	})
}
func (v HouseUnit) ColumnSlice() sequel.ColumnValuer[[]float64] {
	return sequel.Column("slice", v.Slice, func(val []float64) driver.Value {
		return (sqltype.Float64Slice[float64])(val)
	})
}
func (v HouseUnit) ColumnMap() sequel.ColumnValuer[map[string]float64] {
	return sequel.Column("map", v.Map, func(val map[string]float64) driver.Value {
		return encoding.JSONValue(val)
	})
}
