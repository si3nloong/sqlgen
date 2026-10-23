package encoding

import (
	"database/sql"
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

func TestTimeScannerConversions(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	checkScanner(t, func(v *time.Time) sql.Scanner { return TimeScanner(v) }, func(v **time.Time) sql.Scanner { return TimeScanner(v) },
		[]scanCase[time.Time]{{now, now}, {"2026-01-02T03:04:05Z", now}, {[]byte("2026-01-02 03:04:05"), now}, {int64(100), time.Unix(100, 0)}}, []any{"invalid", []byte("invalid"), true})
}
