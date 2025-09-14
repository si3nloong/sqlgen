package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInteger(t *testing.T) {
	t.Run("Scan with nil", func(t *testing.T) {
		var i8 = int8(10)
		v := Int8Scanner[int8](&i8)
		require.NoError(t, v.Scan(nil))
		require.Equal(t, int8(0), i8)
	})

	t.Run("Scan with pointer", func(t *testing.T) {
		var i8 *int8
		val := int8(10)
		i8 = &val
		v := Int8Scanner[int8](&i8)
		require.NoError(t, v.Scan(nil))
		require.Nil(t, i8)
	})

	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("int8", func(t *testing.T) {
			var i8 int8
			v := Int8Scanner[int8](&i8)
			require.NoError(t, v.Scan(int64(8)))
			require.Equal(t, i8, int8(8))
		})

		t.Run("int16", func(t *testing.T) {
			var i16 int16
			v := Int16Scanner[int16](&i16)
			require.NoError(t, v.Scan(int64(-68)))
			require.Equal(t, i16, int16(-68))
		})

		t.Run("int32", func(t *testing.T) {
			var i32 int32
			v := Int32Scanner[int32](&i32)
			require.NoError(t, v.Scan(int64(-128)))
			require.Equal(t, i32, int32(-128))
		})

		t.Run("int64", func(t *testing.T) {
			var i64 int64
			v := Int64Scanner[int64](&i64)
			require.NoError(t, v.Scan(int64(-19_823_028)))
			require.Equal(t, i64, int64(-19_823_028))
		})
	})

	// t.Run("sql.Scanner", func(t *testing.T) {
	// 	var i32 = int32(88)
	// 	v := Integer(&i32)
	// 	require.NoError(t, v.Scan(int64(1580)))
	// })

	// t.Run("driver.Valuer", func(t *testing.T) {
	// 	var i32 = int32(88)
	// 	v := Integer(&i32)
	// 	value, err := v.Value()
	// 	require.NoError(t, err)
	// 	require.Equal(t, int64(88), value)
	// })

	// t.Run("Integer with new(int)", func(t *testing.T) {
	// 	var ptr = new(int)
	// 	v := Integer(ptr)

	// 	t.Run("Value", func(t *testing.T) {
	// 		value, err := v.Value()
	// 		require.NoError(t, err)
	// 		require.Empty(t, value)
	// 	})

	// 	t.Run("Scan", func(t *testing.T) {
	// 		require.NoError(t, v.Scan(nil))
	// 		value, err := v.Value()
	// 		require.NoError(t, err)
	// 		require.Nil(t, value)
	// 	})
	// })
}
