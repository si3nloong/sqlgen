package main

import (
	"database/sql/driver"
	"reflect"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Address) TableName() string {
	return "address"
}
func (Address) Columns() []string {
	return []string{"line_1", "line_2", "country_code"}
}
func (v Address) Values() []any {
	return []any{(string)(v.Line1), (string)(v.Line2), (string)(v.CountryCode)}
}
func (v *Address) Addrs() []any {
	return []any{types.String(&v.Line1), types.String(&v.Line2), types.String(&v.CountryCode)}
}
func (Address) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v Address) InsertOneStmt() (string, []any) {
	return "INSERT INTO address (line_1,line_2,country_code) VALUES (?,?,?);", v.Values()
}
func (v Address) GetLine1() sequel.ColumnValuer[string] {
	return sequel.Column("line_1", v.Line1, func(val string) driver.Value {
		return (string)(val)
	})
}
func (v Address) GetLine2() sequel.ColumnValuer[string] {
	return sequel.Column("line_2", v.Line2, func(val string) driver.Value {
		return (string)(val)
	})
}
func (v Address) GetCountryCode() sequel.ColumnValuer[string] {
	return sequel.Column("country_code", v.CountryCode, func(val string) driver.Value {
		return (string)(val)
	})
}

func (HouseUnit) TableName() string {
	return "house_unit"
}
func (HouseUnit) Columns() []string {
	return []string{"no", "build_time", "address", "kind", "type", "chan", "inner", "arr", "slice", "map"}
}
func (v HouseUnit) Values() []any {
	return []any{(int64)(v.No), (time.Time)(v.BuildTime), types.JSONMarshaler(v.Address), (int64)(v.Kind), (int64)(v.Type), (int64)(v.Chan), types.JSONMarshaler(v.Inner), types.JSONMarshaler(v.Arr), encoding.MarshalFloatList(v.Slice, -1), types.JSONMarshaler(v.Map)}
}
func (v *HouseUnit) Addrs() []any {
	return []any{types.Integer(&v.No), (*time.Time)(&v.BuildTime), types.JSONUnmarshaler(&v.Address), types.Integer(&v.Kind), types.Integer(&v.Type), types.Integer(&v.Chan), types.JSONUnmarshaler(&v.Inner), types.JSONUnmarshaler(&v.Arr), types.Float64Slice(&v.Slice), types.JSONUnmarshaler(&v.Map)}
}
func (HouseUnit) InsertPlaceholders(row int) string {
	return "(?,?,?,?,?,?,?,?,?,?)"
}
func (v HouseUnit) InsertOneStmt() (string, []any) {
	return "INSERT INTO house_unit (no,build_time,address,kind,type,chan,inner,arr,slice,map) VALUES (?,?,?,?,?,?,?,?,?,?);", v.Values()
}
func (v HouseUnit) GetNo() sequel.ColumnValuer[uint] {
	return sequel.Column("no", v.No, func(val uint) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) GetBuildTime() sequel.ColumnValuer[time.Time] {
	return sequel.Column("build_time", v.BuildTime, func(val time.Time) driver.Value {
		return (time.Time)(val)
	})
}
func (v HouseUnit) GetAddress() sequel.ColumnValuer[Address] {
	return sequel.Column("address", v.Address, func(val Address) driver.Value {
		return types.JSONMarshaler(val)
	})
}
func (v HouseUnit) GetKind() sequel.ColumnValuer[reflect.Kind] {
	return sequel.Column("kind", v.Kind, func(val reflect.Kind) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) GetType() sequel.ColumnValuer[HouseUnitType] {
	return sequel.Column("type", v.Type, func(val HouseUnitType) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) GetChan() sequel.ColumnValuer[reflect.ChanDir] {
	return sequel.Column("chan", v.Chan, func(val reflect.ChanDir) driver.Value {
		return (int64)(val)
	})
}
func (v HouseUnit) GetInner() sequel.ColumnValuer[struct{ Flag bool }] {
	return sequel.Column("inner", v.Inner, func(val struct{ Flag bool }) driver.Value {
		return types.JSONMarshaler(val)
	})
}
func (v HouseUnit) GetArr() sequel.ColumnValuer[[2]string] {
	return sequel.Column("arr", v.Arr, func(val [2]string) driver.Value {
		return types.JSONMarshaler(val)
	})
}
func (v HouseUnit) GetSlice() sequel.ColumnValuer[[]float64] {
	return sequel.Column("slice", v.Slice, func(val []float64) driver.Value {
		return encoding.MarshalFloatList(val, -1)
	})
}
func (v HouseUnit) GetMap() sequel.ColumnValuer[map[string]float64] {
	return sequel.Column("map", v.Map, func(val map[string]float64) driver.Value {
		return types.JSONMarshaler(val)
	})
}
