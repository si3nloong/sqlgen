package sql

import (
	"database/sql/driver"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
)

func TestAutoPkLocation(t *testing.T) {
	t.Run("AutoPkLocation has empty values", func(t *testing.T) {
		l := AutoPkLocation{}
		values := l.Values()
		require.Equal(t, 4, len(values))
		require.Nil(t, values[1])
		require.Nil(t, values[2])
		require.Nil(t, values[3])
	})

	t.Run("AutoPkLocation has values", func(t *testing.T) {
		uid, _ := uuid.NewV1()
		l := AutoPkLocation{}
		l.GeoPoint = orb.Point{1, 3}
		l.PtrGeoPoint = &l.GeoPoint
		l.PtrUUID = &uid
		values := l.Values()
		require.Equal(t, 4, len(values))
		_, ok := values[1].(driver.Valuer)
		require.True(t, ok)
		require.NotNil(t, values[2])
		_, ok = values[2].(driver.Valuer)
		require.True(t, ok)
		require.Nil(t, values[3])
	})
}
