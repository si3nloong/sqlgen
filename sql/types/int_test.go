package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInteger(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		var intVal int8
		v := Integer(&intVal)
		require.NoError(t, v.Scan(int64(8)))
		require.Equal(t, int8(8), v.Interface())
	})
}
