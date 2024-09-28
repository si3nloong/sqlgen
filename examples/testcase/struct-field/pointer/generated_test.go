package pointer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointer(t *testing.T) {
	ptr := Ptr{}
	ptr.embeded = new(embeded)
	values := ptr.Values()
	require.Equal(t, 19, len(values))
	require.Equal(t, int64(0), values[0])
	for i := 1; i < 18; i++ {
		require.Nil(t, values[i])
	}
}
