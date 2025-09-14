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
		require.Equal(t, "$1", driver.QuoteVar(1))
		require.Equal(t, "$10", driver.QuoteVar(10))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.Equal(t, `"abc"`, driver.QuoteIdentifier("abc"))
		require.Equal(t, `"abc_def"`, driver.QuoteIdentifier("abc_def"))
	})

	t.Run("QuoteRune", func(t *testing.T) {
		require.Equal(t, rune('"'), driver.QuoteRune())
	})
}
