package encoding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMarshalStringSlice(t *testing.T) {
	t.Run("MarshalStringSlice using []byte", func(t *testing.T) {
		require.Equal(t, `["a","b"]`, MarshalStringSlice([][]byte{[]byte("a"), []byte("b")}))
	})

	t.Run("MarshalStringSlice using string", func(t *testing.T) {
		require.Equal(t, `["a","b"]`, MarshalStringSlice([]string{"a", "b"}))
	})

	t.Run("MarshalStringSlice using custom string", func(t *testing.T) {
		type CustomStr string
		require.Equal(t, `["a","b","z"]`, MarshalStringSlice([]CustomStr{"a", "b", "z"}))
	})
}

func TestMarshalIntSlice(t *testing.T) {
	t.Run("MarshalIntSlice using empty int array", func(t *testing.T) {
		require.Equal(t, `[]`, MarshalIntSlice([]int{}))
	})

	t.Run("MarshalIntSlice using single int array", func(t *testing.T) {
		require.Equal(t, `[-1]`, MarshalIntSlice([]int{-1}))
	})

	t.Run("MarshalIntSlice using int", func(t *testing.T) {
		require.Equal(t, `[-1,-6,11,88,100]`, MarshalIntSlice([]int{-1, -6, 11, 88, 100}))
	})

	t.Run("MarshalIntSlice using uint", func(t *testing.T) {
		require.Equal(t, `[1,5,10]`, MarshalUintSlice([]uint{1, 5, 10}))
	})

	t.Run("MarshalIntSlice using iota", func(t *testing.T) {
		type enum int

		const (
			success enum = iota + 1
			failed
			pending
		)
		require.Equal(t, `[1,2,3]`, MarshalIntSlice([]enum{success, failed, pending}))
	})
}

func TestMarshalBoolSlice(t *testing.T) {
	t.Run("MarshalBoolSlice using empty bool array", func(t *testing.T) {
		require.Equal(t, `[]`, MarshalBoolSlice([]bool{}))
	})

	t.Run("MarshalBoolSlice using bool", func(t *testing.T) {
		require.Equal(t, `[true,false,true]`, MarshalBoolSlice([]bool{true, false, true}))
	})

	t.Run("MarshalBoolSlice using custom bool", func(t *testing.T) {
		type Flag bool
		require.Equal(t, `[false,false,true]`, MarshalBoolSlice([]Flag{false, false, true}))
	})
}

func TestMarshalFloatList(t *testing.T) {
	t.Run("MarshalFloatList using empty float array", func(t *testing.T) {
		require.Equal(t, `[]`, MarshalFloat64List([]float64{}))
	})
	// t.Run("MarshalFloatList using float32", func(t *testing.T) {
	// 	require.Equal(t, `[10.05,-881.33,100.55]`, MarshalFloat64List([]float32{10.05, -881.333, 100.5522}, 2))
	// })

	// t.Run("MarshalFloatList using custom float32", func(t *testing.T) {
	// 	type f32 float32

	// 	require.Equal(t, `[12.4526]`, MarshalFloatList([]f32{12.4526}, 4))
	// })

	t.Run("MarshalFloatList using float64", func(t *testing.T) {
		require.Equal(t, `[10.05,-881.33,100.55]`, MarshalFloat64List([]float64{10.05, -881.333, 100.5522}, 2))
	})

	t.Run("MarshalFloatList using custom float64", func(t *testing.T) {
		type f64 float64
		require.Equal(t, `[12.4526]`, MarshalFloat64List([]f64{12.4526}, 4))
	})
}

func TestMarshalTimeList(t *testing.T) {
	t.Run("MarshalTimeList using empty time array", func(t *testing.T) {
		require.Equal(t, `[]`, MarshalTimeList([]time.Time{}))
	})
}
