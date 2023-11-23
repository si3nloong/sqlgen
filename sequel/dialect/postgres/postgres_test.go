package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgresDriver(t *testing.T) {
	driver := new(postgresDriver)
	t.Run("Driver", func(t *testing.T) {
		require.Equal(t, "postgres", driver.Driver())
	})

	t.Run("Var", func(t *testing.T) {
		require.Equal(t, "$1", driver.Var(1))
		require.Equal(t, "$10", driver.Var(10))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.Equal(t, `"abc"`, driver.Wrap("abc"))
		require.Equal(t, `"abc_def"`, driver.Wrap("abc_def"))
	})

	t.Run("QuoteChar", func(t *testing.T) {
		require.Equal(t, rune('"'), driver.QuoteChar())
	})

	t.Run("VarChar", func(t *testing.T) {
		require.Equal(t, "$", driver.VarChar())
	})
}
