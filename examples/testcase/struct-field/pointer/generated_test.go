package pointer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPointer(t *testing.T) {
	t.Run("Every field is nil value", func(t *testing.T) {
		ptr := Ptr{}
		values := ptr.Values()
		require.Equal(t, 19, len(values))
		require.Equal(t, int64(0), values[0])
		for i := 1; i < 19; i++ {
			require.Nil(t, values[i])
		}
	})

	t.Run("nested values", func(t *testing.T) {
		t.Run("deepNested has value but descendant is nil", func(t *testing.T) {
			ptr := Ptr{}
			ptr.deepNested = &deepNested{}
			values := ptr.Values()
			require.Equal(t, 19, len(values))
			require.Equal(t, int64(0), values[0])
			for i := 1; i < 19; i++ {
				require.Nil(t, values[i])
			}
		})

		t.Run("embedded has value but descendant is nil", func(t *testing.T) {
			ptr := Ptr{}
			ptr.deepNested = &deepNested{
				&embedded{},
			}
			values := ptr.Values()
			require.Equal(t, 19, len(values))
			require.Equal(t, int64(0), values[0])
			for i := 1; i < 19; i++ {
				require.Nil(t, values[i])
			}
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
			require.Equal(t, int64(0), values[0])
			for i := 1; i < 18; i++ {
				require.Nil(t, values[i])
			}
			require.NotNil(t, values[18])
			require.Equal(t, ts.Format(time.RFC3339), (values[18]).(time.Time).Format(time.RFC3339))
		})
	})
}
