package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		var flag bool
		v := Bool(&flag)
		require.NoError(t, v.Scan(true))
		require.Equal(t, true, v.Interface())

		v = Bool(&flag)
		require.NoError(t, v.Scan(false))
		require.Equal(t, false, v.Interface())
	})
}
