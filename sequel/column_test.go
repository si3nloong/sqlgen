package sequel

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValueColumn(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value sql.NullString
		want  any
	}{
		{"valid", sql.NullString{String: "Alice", Valid: true}, "Alice"},
		{"null", sql.NullString{}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := ValueColumn("name", tc.value)
			require.Equal(t, "name", c.ColumnName())
			require.Equal(t, tc.value, c.Value())
			value, err := c.Value().Value()
			require.NoError(t, err)
			require.Equal(t, tc.want, value)
		})
	}
}

func TestColumn(t *testing.T) {
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

	t.Run("PrimitiveColumn", func(t *testing.T) {
		bc := PrimitiveColumn(columnName, text)
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

func TestSQLColumn(t *testing.T) {
	c := SQLColumn("name", "Alice", func(p string) string { return "LOWER(" + p + ")" }, func(v string) any { return "converted:" + v })
	require.Equal(t, "name", c.ColumnName())
	require.Equal(t, "Alice", c.Value())
	require.Equal(t, "LOWER($1)", c.SQLColumn("$1"))
	require.Equal(t, "converted:Bob", c.Convert("Bob"))
}

func TestOrderByDescending(t *testing.T) {
	require.True(t, OrderByColumn("id", false).Desc())
	require.False(t, OrderByColumn("id", true).Desc())
}
