package pointer

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPointer(t *testing.T) {
	t.Run("Columns and Addrs count must tally", func(t *testing.T) {
		ptr := Ptr{}
		require.Equal(t, len(ptr.Columns()), len(ptr.Addrs()))
	})
	t.Run("Values", func(t *testing.T) {
		t.Run("Every field is nil value", func(t *testing.T) {
			ptr := Ptr{}
			values := ptr.Values()
			require.Equal(t, 19, len(values))
			// require.Equal(t, int64(0), values[0])
			// for i := 1; i < 19; i++ {
			// 	require.Nil(t, values[i])
			// }
		})

		t.Run("nested values", func(t *testing.T) {
			t.Run("deepNested has value but descendant is nil", func(t *testing.T) {
				ptr := Ptr{}
				ptr.deepNested = &deepNested{}
				require.Nil(t, ptr.GetEmbeddedTime())
				require.Nil(t, ptr.GetAnyTime())
				// values := ptr.Values()
				// require.Equal(t, 19, len(values))
				// require.Equal(t, int64(0), values[0])
				// for i := 1; i < 19; i++ {
				// 	require.Nil(t, values[i])
				// }
			})

			t.Run("embedded has value but descendant is nil", func(t *testing.T) {
				ptr := Ptr{}
				ptr.deepNested = &deepNested{
					&embedded{},
				}
				require.Nil(t, ptr.GetEmbeddedTime())
				require.Zero(t, ptr.GetAnyTime())
				// values := ptr.Values()
				// require.Equal(t, 19, len(values))
				// require.Equal(t, int64(0), values[0])
				// for i := 1; i < 19; i++ {
				// 	require.Nil(t, values[i])
				// }
			})

			t.Run("EmbeddedTime has value", func(t *testing.T) {
				ptr := Ptr{}
				ts := time.Now()
				ptr.deepNested = &deepNested{
					&embedded{
						EmbeddedTime: &ts,
					},
				}
				values := ptr.Values()
				require.Equal(t, 19, len(values))
				require.NotNil(t, ptr.GetEmbeddedTime())
				require.Equal(t, ts.Format(time.RFC3339), ptr.GetEmbeddedTime().(time.Time).Format(time.RFC3339))
			})
		})
	})

	t.Run("Addrs", func(t *testing.T) {
		ptr := Ptr{}
		ptr.Int = new(int)
		ptr.Int8 = new(int8)
		ptr.Int16 = new(int16)
		ptr.Int32 = new(int32)
		ptr.Uint = new(uint)
		ptr.Uint8 = new(uint8)
		ptr.Uint16 = new(uint16)
		ptr.Uint64 = new(uint64)

		addrs := ptr.Addrs()
		require.NotNil(t, ptr.Int)
		require.NoError(t, addrs[4].(sql.Scanner).Scan(nil))
		require.Nil(t, ptr.Int)

		require.NotNil(t, ptr.Int8)
		require.NoError(t, addrs[5].(sql.Scanner).Scan(nil))
		require.Nil(t, ptr.Int8)

		require.NotNil(t, ptr.Int16)
		require.NoError(t, addrs[6].(sql.Scanner).Scan(nil))
		require.Nil(t, ptr.Int16)

		require.NotNil(t, ptr.EmbeddedTime)
		require.NoError(t, addrs[18].(sql.Scanner).Scan(nil))
		require.Nil(t, ptr.EmbeddedTime)

		require.Zero(t, ptr.GetAnyTime())
		require.Zero(t, ptr.AnyTime)
	})
}
