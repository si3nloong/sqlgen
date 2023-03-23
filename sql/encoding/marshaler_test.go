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

	t.Run("MarshalStringList using custom string", func(t *testing.T) {
		type CustomStr string
		require.Equal(t, `["a","b","z"]`, MarshalStringList([]CustomStr{"a", "b", "z"}))
	})
}

func TestMarshalIntList(t *testing.T) {
	t.Run("MarshalIntList using int", func(t *testing.T) {
		require.Equal(t, `[-1,-6,11,88,100]`, MarshalIntList([]int{-1, -6, 11, 88, 100}))
	})

	t.Run("MarshalIntList using uint", func(t *testing.T) {
		require.Equal(t, `[1,5,10]`, MarshalIntList([]uint{1, 5, 10}))
	})

	t.Run("MarshalIntList using iota", func(t *testing.T) {
		type enum int

		const (
			success enum = iota + 1
			failed
			pending
		)
		require.Equal(t, `[1,2,3]`, MarshalIntList([]enum{success, failed, pending}))
	})
}

func TestMarshalBoolList(t *testing.T) {
	t.Run("MarshalBoolList using bool", func(t *testing.T) {
		require.Equal(t, `[true,false,true]`, MarshalBoolList([]bool{true, false, true}))
	})

	t.Run("MarshalBoolList using custom bool", func(t *testing.T) {
		type Flag bool

		require.Equal(t, `[false,false,true]`, MarshalBoolList([]Flag{false, false, true}))
	})
}

func TestMarshalFloatList(t *testing.T) {
	t.Run("MarshalFloatList using float32", func(t *testing.T) {
		require.Equal(t, `[10.05,-881.33,100.55]`, MarshalFloatList([]float32{10.05, -881.333, 100.5522}, 2))
	})

	t.Run("MarshalFloatList using float64", func(t *testing.T) {
		require.Equal(t, `[10.05,-881.33,100.55]`, MarshalFloatList([]float64{10.05, -881.333, 100.5522}, 2))
	})

	t.Run("MarshalFloatList using custom float32", func(t *testing.T) {
		type f32 float32

		require.Equal(t, `[12.4526]`, MarshalFloatList([]f32{12.4526}, 4))
	})
}

func TestMarshalTimeList(t *testing.T) {

}
