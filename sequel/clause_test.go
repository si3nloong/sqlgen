package sequel

import (
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

	cv := Column(columnName, Str(value), func(s Str) any {
		return string(s)
	})
	require.Equal(t, columnName, cv.ColumnName())
	require.Equal(t, value, cv.Value())
	require.Equal(t, text, cv.Convert(Str(text)))
}
