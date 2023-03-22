package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalStringList(t *testing.T) {
	t.Run("MarshalStringList using []byte", func(t *testing.T) {
		require.Equal(t, `["a","b"]`, MarshalStringList([][]byte{[]byte("a"), []byte("b")}))
	})

	t.Run("MarshalStringList using string", func(t *testing.T) {
		require.Equal(t, `["a","b"]`, MarshalStringList([]string{"a", "b"}))
	})
}

func TestMarshalIntList(t *testing.T) {
	t.Run("MarshalIntList using int", func(t *testing.T) {
		require.Equal(t, `[-1,-6,11,88,100]`, MarshalIntList([]int{-1, -6, 11, 88, 100}))
	})

	t.Run("MarshalIntList using uint", func(t *testing.T) {
		require.Equal(t, `[1,5,10]`, MarshalIntList([]uint{1, 5, 10}))
	})
}
