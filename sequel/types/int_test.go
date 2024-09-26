package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInteger(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("int8", func(t *testing.T) {
			var i8 int8
			v := Integer(&i8)
			require.NoError(t, v.Scan(int64(8)))
			require.Equal(t, int8(8), v.Interface())
		})

		t.Run("int16", func(t *testing.T) {
			var i16 int16
			v := Integer(&i16)
			require.NoError(t, v.Scan(-68))
			require.Equal(t, int16(-68), v.Interface())
		})

		t.Run("int32", func(t *testing.T) {
			var i32 int32
			v := Integer(&i32)
			require.NoError(t, v.Scan(-128))
			require.Equal(t, int32(-128), v.Interface())
		})
	})

	t.Run("sql.Scanner", func(t *testing.T) {
		var i32 = int32(88)
		v := Integer(&i32)
		require.NoError(t, v.Scan(int64(1580)))
		require.Equal(t, int32(1580), v.Interface())
	})

	t.Run("driver.Valuer", func(t *testing.T) {
		var i32 = int32(88)
		v := Integer(&i32)
		value, err := v.Value()
		require.NoError(t, err)
		require.Equal(t, int64(88), value)
	})

	t.Run("Integer with new(int)", func(t *testing.T) {
		var ptr = new(int)
		v := Integer(ptr)

		t.Run("Value", func(t *testing.T) {
			value, err := v.Value()
			require.NoError(t, err)
			require.Empty(t, value)
		})

		t.Run("Scan", func(t *testing.T) {
			require.NoError(t, v.Scan(nil))
			value, err := v.Value()
			require.NoError(t, err)
			require.Nil(t, value)
		})
	})
}
