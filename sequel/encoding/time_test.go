package encoding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {

	type Alias = time.Time

	var ts Alias
	require.NoError(t, TimeScanner(&ts).Scan(`2024-01-01`))
	require.Equal(t, `2024-01-01`, ts.Format("2006-01-02"))

	var ptrTs *Alias
	require.NoError(t, TimeScanner(&ptrTs).Scan(nil))
	require.Nil(t, ptrTs)

	t.Run("Initialized of Time", func(t *testing.T) {
		var ts time.Time
		require.NoError(t, TimeScanner(&ts).Scan(nil))
		require.NotNil(t, ts)
		require.Zero(t, ts)

		t.Run("Scan nil value", func(t *testing.T) {
			var ts = new(time.Time)
			require.NoError(t, TimeScanner(&ts).Scan(nil))
			require.Nil(t, ts)

			var ptrTs = new(time.Time)
			require.NoError(t, TimeScanner(&ptrTs).Scan(nil))
			require.Nil(t, ptrTs)
		})

		t.Run("Scan value", func(t *testing.T) {
			var ts = new(time.Time)
			require.NoError(t, TimeScanner(&ts).Scan(`2008-01-31`))
			require.Equal(t, `2008-01-31`, ts.Format(`2006-01-02`))
		})
	})

	t.Run("Initialized Pointer of Time", func(t *testing.T) {
		t.Run("Scan nil value", func(t *testing.T) {
			var ts = new(time.Time)
			require.NoError(t, TimeScanner(&ts).Scan(nil))
			require.Nil(t, ts)
		})

		t.Run("Scan value", func(t *testing.T) {
			var ts = new(time.Time)
			require.NoError(t, TimeScanner(&ts).Scan(`2024-01-01`))
			require.Equal(t, `2024-01-01`, ts.Format("2006-01-02"))
		})
	})

}
