package sql

import (
	"cloud.google.com/go/civil"
	"github.com/gofrs/uuid/v5"
	"github.com/paulmach/orb"
)

type Location struct {
	ID       uint64 `sql:",pk"`
	GeoPoint orb.Point
	UUID     uuid.UUID
}

type AutoPkLocation struct {
	ID          uint64 `sql:",pk,auto_increment"`
	GeoPoint    orb.Point
	PtrGeoPoint *orb.Point
	PtrUUID     *uuid.UUID
	PtrDate     *civil.Date
}
