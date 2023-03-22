package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAssertAs(t *testing.T) {
	t.Run("assertAs should returns something", func(t *testing.T) {
		str := ""
		v := assertAs[string](&str)
		require.NotNil(t, v)
	})

	t.Run("assertAs should returns nil", func(t *testing.T) {
		str := ""
		v := assertAs[string](str)
		require.Nil(t, v)
	})
}
