package aliasstruct

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/require"
)

func TestA(t *testing.T) {
	a := A{}

	t.Run("Values (TextMarshaler)", func(t *testing.T) {
		t.Run("civil.Date", func(t *testing.T) {
			val, err := a.DateValue().(driver.Valuer).Value()
			require.NoError(t, err)
			require.Equal(t, []byte(`0000-00-00`), val)
		})

		t.Run("civil.Time", func(t *testing.T) {
			val, err := a.TimeValue().(driver.Valuer).Value()
			require.NoError(t, err)
			require.Equal(t, []byte(`00:00:00`), val)
		})
	})

	t.Run("Addrs (TextUnmarshaler)", func(t *testing.T) {
		addrs := a.Addrs()

		t.Run("civil.Date", func(t *testing.T) {
			d, _ := civil.ParseDate(`2004-03-14`)
			require.NoError(t, addrs[0].(sql.Scanner).Scan(`2004-03-14`))
			require.Equal(t, d, a.Date)
		})

		t.Run("civil.Time", func(t *testing.T) {
			ts, _ := civil.ParseTime(`15:03:09`)
			require.NoError(t, addrs[1].(sql.Scanner).Scan(`15:03:09`))
			require.Equal(t, ts, a.Time)
		})
	})
}
