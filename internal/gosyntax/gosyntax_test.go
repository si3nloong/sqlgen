package gosyntax

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPtrOf(t *testing.T) {
	require.Equal(t, "abc", *PtrOf("abc"))
	require.Equal(t, 8, *PtrOf(8))
	require.Equal(t, uint(88), *PtrOf(uint(88)))
	require.NotNil(t, PtrOf(""))
}
