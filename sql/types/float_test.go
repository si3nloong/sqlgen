package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			var f32 float32
			v := Float(&f32)
			require.NoError(t, v.Scan(float64(81.20022)))
			require.Equal(t, float32(81.20022), v.Interface())
		})
	})
}
