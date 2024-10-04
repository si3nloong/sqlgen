package valuer

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("Columns", func(t *testing.T) {
		b := B{}
		require.ElementsMatch(t, []string{"id", "value", "ptr_value", "n"}, b.Columns())
	})

	t.Run("Values with empty values", func(t *testing.T) {
		b := B{}
		values := b.Values()
		require.Equal(t, 4, len(values))
		require.Equal(t, int64(0), values[0])
		_, ok := values[1].(driver.Valuer)
		require.True(t, ok)
		require.Nil(t, values[2])
		require.Equal(t, "", values[3])
	})

	t.Run("Values with values", func(t *testing.T) {
		b := B{}
		b.ID = 100
		b.Value = anyType{}
		b.PtrValue = &anyType{ptr: true}
		b.N = "Hello WORLD!!"

		values := b.Values()
		require.Equal(t, 4, len(values))
		require.Equal(t, int64(100), values[0])
		v, ok := values[1].(driver.Valuer)
		require.True(t, ok)
		val, err := v.Value()
		require.NoError(t, err)
		require.Equal(t, "any", val)
		require.NotNil(t, values[2])
		v, ok = values[2].(driver.Valuer)
		require.True(t, ok)
		val, err = v.Value()
		require.NoError(t, err)
		require.Equal(t, "ptr", val)
		require.Equal(t, "Hello WORLD!!", values[3])
	})

	t.Run("Addrs shouldn't be nil", func(t *testing.T) {
		b := B{}
		for _, addr := range b.Addrs() {
			require.NotNil(t, addr)
		}
	})
}
