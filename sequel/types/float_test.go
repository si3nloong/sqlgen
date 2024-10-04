package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			var f32 float32
			v := Float32(&f32)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.Equal(t, float32(81.20022), v.Interface())
		})

		t.Run("~float32", func(t *testing.T) {
			type F32 float32
			var f F32
			v := Float32(&f)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.Equal(t, F32(81.20022), v.Interface())
		})

		t.Run("float64", func(t *testing.T) {
			var f64 float64
			v := Float64(&f64)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.Equal(t, float64(81.20022), v.Interface())
		})
	})
}
