package sequel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResult(t *testing.T) {
	t.Run("RowsAffected", func(t *testing.T) {
		result := result{rowsAffected: 100}
		affected, err := result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(100), affected)
	})
}

func TestEmptyResult(t *testing.T) {
	result := new(EmptyResult)

	n, err := result.LastInsertId()
	require.NoError(t, err)
	require.Equal(t, int64(0), n)

	affected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(0), affected)
}

func TestRowsAffectedResultContract(t *testing.T) {
	for _, count := range []int64{0, 1, 100} {
		r := NewRowsAffectedResult(count)
		got, err := r.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, count, got)
		id, err := r.LastInsertId()
		require.NoError(t, err)
		require.Zero(t, id)
	}
	r := result{lastID: 42}
	id, err := r.LastInsertId()
	require.NoError(t, err)
	require.Equal(t, int64(42), id)
}
