package sql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStmt(t *testing.T) {
	t.Run("WriteQuery", func(t *testing.T) {
		stmt := AcquireStmt()
		stmt.WriteQuery(`abcde!`, int32(8))
		require.Equal(t, `abcde!`, stmt.Query())
		require.ElementsMatch(t, []any{int32(8)}, stmt.Args())
		ReleaseStmt(stmt)

		stmt = AcquireStmt()
		require.Empty(t, stmt.Query())
		require.Empty(t, stmt.Args())
		ReleaseStmt(stmt)
	})

	t.Run("Reset", func(t *testing.T) {
		stmt := AcquireStmt()
		defer ReleaseStmt(stmt)
		stmt.WriteQuery(`abcde!`, int32(8))
		require.Equal(t, `abcde!`, stmt.Query())
		require.ElementsMatch(t, []any{int32(8)}, stmt.Args())

		stmt.Reset()
		require.Empty(t, stmt.Query())
		require.Empty(t, stmt.Args())
	})
}
