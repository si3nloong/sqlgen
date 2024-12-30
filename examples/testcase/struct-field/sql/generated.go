package sql

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (AutoPkLocation) TableName() string {
	return "auto_pk_location"
}
func (AutoPkLocation) HasPK()      {}
func (AutoPkLocation) IsAutoIncr() {}
func (v *AutoPkLocation) ScanAutoIncr(val int64) error {
	v.ID = uint64(val)
	return nil
}
func (v AutoPkLocation) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (AutoPkLocation) Columns() []string {
	return []string{"id", "geo_point", "ptr_geo_point", "ptr_uuid", "ptr_date"} // 5
}
func (v AutoPkLocation) Values() []any {
	return []any{
		ewkb.Value(v.GeoPoint, 4326), // 1 - geo_point
		v.PtrGeoPointValue(),         // 2 - ptr_geo_point
		v.PtrUUIDValue(),             // 3 - ptr_uuid
		v.PtrDateValue(),             // 4 - ptr_date
	}
}
func (v *AutoPkLocation) Addrs() []any {
	if v.PtrGeoPoint == nil {
		v.PtrGeoPoint = new(orb.Point)
	}
	if v.PtrUUID == nil {
		v.PtrUUID = new(uuid.UUID)
	}
	if v.PtrDate == nil {
		v.PtrDate = new(civil.Date)
	}
	return []any{
		encoding.Uint64Scanner[uint64](&v.ID),        // 0 - id
		ewkb.Scanner(&v.GeoPoint),                    // 1 - geo_point
		encoding.JSONScanner(&v.PtrGeoPoint),         // 2 - ptr_geo_point
		encoding.PtrScanner(&v.PtrUUID),              // 3 - ptr_uuid
		encoding.TextScanner[civil.Date](&v.PtrDate), // 4 - ptr_date
	}
}
func (AutoPkLocation) InsertColumns() []string {
	return []string{"geo_point", "ptr_geo_point", "ptr_uuid", "ptr_date"} // 4
}
func (AutoPkLocation) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v AutoPkLocation) InsertOneStmt() (string, []any) {
	return "INSERT INTO auto_pk_location (geo_point,ptr_geo_point,ptr_uuid,ptr_date) VALUES (?,?,?,?);", []any{ewkb.Value(v.GeoPoint, 4326), v.PtrGeoPointValue(), v.PtrUUIDValue(), v.PtrDateValue()}
}
func (v AutoPkLocation) FindOneByPKStmt() (string, []any) {
	return "SELECT id,geo_point,ptr_geo_point,ptr_uuid,ptr_date FROM auto_pk_location WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v AutoPkLocation) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE auto_pk_location SET geo_point = ?,ptr_geo_point = ?,ptr_uuid = ?,ptr_date = ? WHERE id = ?;", []any{ewkb.Value(v.GeoPoint, 4326), v.PtrGeoPointValue(), v.PtrUUIDValue(), v.PtrDateValue(), (int64)(v.ID)}
}
func (v AutoPkLocation) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v AutoPkLocation) GeoPointValue() driver.Value {
	return ewkb.Value(v.GeoPoint, 4326)
}
func (v AutoPkLocation) PtrGeoPointValue() driver.Value {
	if v.PtrGeoPoint != nil {
		return ewkb.Value(*v.PtrGeoPoint, 4326)
	}
	return nil
}
func (v AutoPkLocation) PtrUUIDValue() driver.Value {
	if v.PtrUUID != nil {
		return *v.PtrUUID
	}
	return nil
}
func (v AutoPkLocation) PtrDateValue() driver.Value {
	if v.PtrDate != nil {
		return encoding.TextValue(*v.PtrDate)
	}
	return nil
}
func (v AutoPkLocation) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value {
		return (int64)(val)
	})
}
func (v AutoPkLocation) GetGeoPoint() sequel.ColumnValuer[orb.Point] {
	return sequel.Column("geo_point", v.GeoPoint, func(val orb.Point) driver.Value {
		return ewkb.Value(val, 4326)
	})
}
func (v AutoPkLocation) GetPtrGeoPoint() sequel.ColumnValuer[*orb.Point] {
	return sequel.Column("ptr_geo_point", v.PtrGeoPoint, func(val *orb.Point) driver.Value {
		if val != nil {
			return ewkb.Value(*val, 4326)
		}
		return nil
	})
}
func (v AutoPkLocation) GetPtrUUID() sequel.ColumnValuer[*uuid.UUID] {
	return sequel.Column("ptr_uuid", v.PtrUUID, func(val *uuid.UUID) driver.Value {
		if val != nil {
			return *val
		}
		return nil
	})
}
func (v AutoPkLocation) GetPtrDate() sequel.ColumnValuer[*civil.Date] {
	return sequel.Column("ptr_date", v.PtrDate, func(val *civil.Date) driver.Value {
		if val != nil {
			return encoding.TextValue(*val)
		}
		return nil
	})
}

func (Location) TableName() string {
	return "location"
}
func (Location) HasPK() {}
func (v Location) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (Location) Columns() []string {
	return []string{"id", "geo_point", "uuid"} // 3
}
func (v Location) Values() []any {
	return []any{
		(int64)(v.ID),                // 0 - id
		ewkb.Value(v.GeoPoint, 4326), // 1 - geo_point
		v.UUID,                       // 2 - uuid
	}
}
func (v *Location) Addrs() []any {
	return []any{
		encoding.Uint64Scanner[uint64](&v.ID), // 0 - id
		ewkb.Scanner(&v.GeoPoint),             // 1 - geo_point
		&v.UUID,                               // 2 - uuid
	}
}
func (Location) InsertPlaceholders(row int) string {
	return "(?,?,?)" // 3
}
func (v Location) InsertOneStmt() (string, []any) {
	return "INSERT INTO location (id,geo_point,uuid) VALUES (?,?,?);", v.Values()
}
func (v Location) FindOneByPKStmt() (string, []any) {
	return "SELECT id,geo_point,uuid FROM location WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v Location) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE location SET geo_point = ?,uuid = ? WHERE id = ?;", []any{ewkb.Value(v.GeoPoint, 4326), v.UUID, (int64)(v.ID)}
}
func (v Location) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v Location) GeoPointValue() driver.Value {
	return ewkb.Value(v.GeoPoint, 4326)
}
func (v Location) UUIDValue() driver.Value {
	return v.UUID
}
func (v Location) GetID() sequel.ColumnValuer[uint64] {
	return sequel.Column("id", v.ID, func(val uint64) driver.Value {
		return (int64)(val)
	})
}
func (v Location) GetGeoPoint() sequel.ColumnValuer[orb.Point] {
	return sequel.Column("geo_point", v.GeoPoint, func(val orb.Point) driver.Value {
		return ewkb.Value(val, 4326)
	})
}
func (v Location) GetUUID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("uuid", v.UUID, func(val uuid.UUID) driver.Value {
		return val
	})
}
