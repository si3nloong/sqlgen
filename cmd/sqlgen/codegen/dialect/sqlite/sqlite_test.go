package sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqliteDriver(t *testing.T) {
	driver := new(sqliteDriver)
	t.Run("Driver", func(t *testing.T) {
		require.Equal(t, "sqlite", driver.Driver())
	})

	t.Run("Var", func(t *testing.T) {
		require.Equal(t, "?", driver.QuoteVar(0))
		require.Equal(t, "?", driver.QuoteVar(10))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.Equal(t, `"abc"`, driver.QuoteIdentifier("abc"))
		require.Equal(t, `"abc_def"`, driver.QuoteIdentifier("abc_def"))
	})

	t.Run("QuoteRune", func(t *testing.T) {
		require.Equal(t, rune('"'), driver.QuoteRune())
	})
}
