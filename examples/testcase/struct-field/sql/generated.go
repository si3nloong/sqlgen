package sql

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (AutoPkLocation) TableName() string {
	return "auto_pk_location"
}
func (AutoPkLocation) HasPK()      {}
func (AutoPkLocation) IsAutoIncr() {}
func (v AutoPkLocation) PK() (string, int, any) {
	return "id", 0, int64(v.ID)
}
func (AutoPkLocation) Columns() []string {
	return []string{"id", "geo_point", "ptr_geo", "ptr_date"}
}
func (AutoPkLocation) SQLColumns() []string {
	return []string{"id", "ST_AsBinary(geo_point, 4326)", "ptr_geo", "ptr_date"}
}
func (v AutoPkLocation) Values() []any {
	return []any{int64(v.ID), ewkb.Value(v.GeoPoint, 4326), types.JSONMarshaler(v.PtrGeo), types.TextMarshaler(v.PtrDate)}
}
func (v *AutoPkLocation) Addrs() []any {
	return []any{types.Integer(&v.ID), ewkb.Scanner(&v.GeoPoint), types.JSONUnmarshaler(&v.PtrGeo), types.TextUnmarshaler(&v.PtrDate)}
}
func (AutoPkLocation) InsertPlaceholders(row int) string {
	return "(?,?,?)"
}
func (v AutoPkLocation) InsertOneStmt() (string, []any) {
	return "INSERT INTO auto_pk_location (geo_point,ptr_geo,ptr_date) VALUES (ST_GeomFromEWKB(?),?,?);", []any{ewkb.Value(v.GeoPoint, 4326), types.JSONMarshaler(v.PtrGeo), types.TextMarshaler(v.PtrDate)}
}
func (v AutoPkLocation) FindOneByPKStmt() (string, []any) {
	return "SELECT id,ST_AsBinary(geo_point, 4326),ptr_geo,ptr_date FROM auto_pk_location WHERE id = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v AutoPkLocation) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE auto_pk_location SET geo_point = ST_GeomFromEWKB(?),ptr_geo = ?,ptr_date = ? WHERE id = ?;", []any{ewkb.Value(v.GeoPoint, 4326), types.JSONMarshaler(v.PtrGeo), types.TextMarshaler(v.PtrDate), int64(v.ID)}
}
func (v AutoPkLocation) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value { return int64(val) })
}
func (v AutoPkLocation) GetGeoPoint() sequel.SQLColumnValuer[orb.Point] {
	return sequel.SQLColumn("geo_point", v.GeoPoint, func(placeholder string) string { return "ST_GeomFromEWKB(" + placeholder + ")" }, func(val orb.Point) driver.Value { return ewkb.Value(val, 4326) })
}
func (v AutoPkLocation) GetPtrGeo() sequel.ColumnValuer[*orb.Point] {
	return sequel.Column("ptr_geo", v.PtrGeo, func(val *orb.Point) driver.Value { return types.JSONMarshaler(val) })
}
func (v AutoPkLocation) GetPtrDate() sequel.ColumnValuer[*civil.Date] {
	return sequel.Column("ptr_date", v.PtrDate, func(val *civil.Date) driver.Value { return types.TextMarshaler(val) })
}

func (Location) TableName() string {
	return "location"
}
func (Location) HasPK() {}
func (v Location) PK() (string, int, any) {
	return "id", 0, int64(v.ID)
}
func (Location) Columns() []string {
	return []string{"id", "geo_point"}
}
func (Location) SQLColumns() []string {
	return []string{"id", "ST_AsBinary(geo_point, 4326)"}
}
func (v Location) Values() []any {
	return []any{int64(v.ID), ewkb.Value(v.GeoPoint, 4326)}
}
func (v *Location) Addrs() []any {
	return []any{types.Integer(&v.ID), ewkb.Scanner(&v.GeoPoint)}
}
func (Location) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v Location) InsertOneStmt() (string, []any) {
	return "INSERT INTO location (id,geo_point) VALUES (?,ST_GeomFromEWKB(?));", v.Values()
}
func (v Location) FindOneByPKStmt() (string, []any) {
	return "SELECT id,ST_AsBinary(geo_point, 4326) FROM location WHERE id = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v Location) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE location SET geo_point = ST_GeomFromEWKB(?) WHERE id = ?;", []any{ewkb.Value(v.GeoPoint, 4326), int64(v.ID)}
}
func (v Location) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value { return int64(val) })
}
func (v Location) GetGeoPoint() sequel.SQLColumnValuer[orb.Point] {
	return sequel.SQLColumn("geo_point", v.GeoPoint, func(placeholder string) string { return "ST_GeomFromEWKB(" + placeholder + ")" }, func(val orb.Point) driver.Value { return ewkb.Value(val, 4326) })
}
