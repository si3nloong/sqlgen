package sequel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClause(t *testing.T) {
	type Str string

	var (
		columnName = "Address"
		value      = `Jalan Petaling Jaya`
		text       = "Hello World"
	)

	t.Run("Column", func(t *testing.T) {
		cv := Column(columnName, Str(value), func(s Str) any {
			return string(s)
		})
		require.Equal(t, columnName, cv.ColumnName())
		require.Equal(t, Str(value), cv.Value())
		require.Equal(t, text, cv.Convert(Str(text)))
	})

	t.Run("BasicColumn", func(t *testing.T) {
		bc := BasicColumn(columnName, text)
		require.Equal(t, columnName, bc.ColumnName())
		require.Equal(t, text, bc.Value())
	})

	t.Run("OrderByColumn", func(t *testing.T) {
		ob := OrderByColumn(columnName, true)
		require.Equal(t, columnName, ob.ColumnName())
		require.True(t, ob.Asc())

		ob = OrderByColumn(columnName, false)
		require.Equal(t, columnName, ob.ColumnName())
		require.False(t, ob.Asc())
	})
}
