package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolSlice(t *testing.T) {
	t.Run("primitive bool", func(t *testing.T) {
		var bList = []bool{true, false, true}
		var v = BoolSlice(&bList)
		require.NoError(t, v.Scan(nullBytes))
		value, err := v.Value()
		require.NoError(t, err)
		require.Equal(t, nullBytes, value)
	})

	t.Run("custom bool type", func(t *testing.T) {

	})
}
