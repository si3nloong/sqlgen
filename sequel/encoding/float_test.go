package encoding

import (
	"database/sql"
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

func TestFloat32ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *float32) sql.Scanner { return Float32Scanner[float32](v) }, func(v **float32) sql.Scanner { return Float32Scanner[float32](v) },
		[]scanCase[float32]{{float32(1.5), 1.5}, {float64(-2.25), -2.25}, {int64(42), 42}, {"1.25", 1.25}, {[]byte("2.5"), 2.5}}, []any{true, "invalid", []byte("invalid")})
	var f float32
	require.Error(t, Float32Scanner[float32](&f, true).Scan("1.25"))
	require.NoError(t, Float32Scanner[float32](&f, true).Scan(1.25))
	require.Equal(t, float32(1.25), f)
}

func TestFloat64ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *float64) sql.Scanner { return Float64Scanner[float64](v) }, func(v **float64) sql.Scanner { return Float64Scanner[float64](v) },
		[]scanCase[float64]{{float32(1.5), 1.5}, {float64(-2.25), -2.25}, {int64(42), 42}, {"1.25", 1.25}, {[]byte("2.5"), 2.5}}, []any{true, "invalid", []byte("invalid")})
	var f float64
	require.Error(t, Float64Scanner[float64](&f, true).Scan("1.25"))
	require.NoError(t, Float64Scanner[float64](&f, true).Scan(1.25))
	require.Equal(t, float64(1.25), f)
}
