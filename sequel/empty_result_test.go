package sequel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	result := new(EmptyResult)

	n, err := result.LastInsertId()
	require.NoError(t, err)
	require.Equal(t, int64(0), n)

	affected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(0), affected)
}
