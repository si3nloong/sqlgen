package sequel

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	type Str string

	var (
		columnName = "Address"
		value      = `Jalan Petaling Jaya`
		text       = "Hello World"
	)

	cv := Column[Str](columnName, Str(value), func(s Str) driver.Value {
		return string(s)
	})
	require.Equal(t, columnName, cv.ColumnName())
	require.Equal(t, value, cv.Value())
	require.Equal(t, text, cv.Convert(Str(text)))
}
