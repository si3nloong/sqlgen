package sql

import "github.com/paulmach/orb"

type Location struct {
	ID       uint64 `sql:",pk"`
	GeoPoint orb.Point
}

type AutoPkLocation struct {
	ID       uint64 `sql:",pk,auto_increment"`
	GeoPoint orb.Point
}
