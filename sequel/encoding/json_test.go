package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type A struct {
	B *B
}

type B struct {
	ID int64
}

func TestJSON(t *testing.T) {
	a := A{}
	if a.B == nil {
		a.B = new(B)
	}
	require.NoError(t, JSONScanner(&a.B).Scan(nil))
	require.Nil(t, a.B)
}
