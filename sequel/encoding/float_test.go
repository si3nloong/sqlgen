package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat32Scanner(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			var f32 float32
			v := Float32Scanner[float32](&f32)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.GreaterOrEqual(t, f32, float32(0.0))
		})

		t.Run("~float32", func(t *testing.T) {
			type F32 float32
			var f F32
			v := Float32Scanner[F32](&f)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.GreaterOrEqual(t, f, float32(0.0))
		})
	})
}

func TestFloat64Scanner(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("float64", func(t *testing.T) {
			var f64 float64
			v := Float64Scanner[float64](&f64)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.GreaterOrEqual(t, f64, float64(0.0))
		})

		t.Run("~float64", func(t *testing.T) {
			type F64 float64
			var f F64
			v := Float64Scanner[F64](&f)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.GreaterOrEqual(t, f, float64(0.0))
		})
	})
}
