package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysqlDriver(t *testing.T) {
	driver := new(mysqlDriver)
	t.Run("Driver", func(t *testing.T) {
		require.Equal(t, "mysql", driver.Driver())
	})

	t.Run("Var", func(t *testing.T) {
		require.Equal(t, "?", driver.Var(0))
		require.Equal(t, "?", driver.Var(10))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.Equal(t, "`abc`", driver.Wrap("abc"))
		require.Equal(t, "`abc_def`", driver.Wrap("abc_def"))
	})

}
