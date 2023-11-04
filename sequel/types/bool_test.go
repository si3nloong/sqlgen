package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	t.Run("Scan with bool", func(t *testing.T) {
		var flag bool
		v := Bool(&flag)
		require.NoError(t, v.Scan(true))
		require.True(t, flag)
		require.True(t, v.Interface())

		v = Bool(&flag)
		require.NoError(t, v.Scan(false))
		require.False(t, flag)
		require.False(t, v.Interface())
	})

	t.Run("Scan with []byte", func(t *testing.T) {
		var flag bool
		v := Bool(&flag)
		require.NoError(t, v.Scan([]byte(`true`)))
		require.True(t, flag)
		require.True(t, v.Interface())

		v = Bool(&flag)
		require.NoError(t, v.Scan([]byte(`false`)))
		require.False(t, flag)
		require.False(t, v.Interface())
	})

	t.Run("Scan with int64", func(t *testing.T) {
		var flag bool
		v := Bool(&flag)
		require.NoError(t, v.Scan(int64(1)))
		require.True(t, flag)
		require.True(t, v.Interface())

		v = Bool(&flag)
		require.NoError(t, v.Scan(int64(0)))
		require.False(t, flag)
		require.False(t, v.Interface())
	})
}
