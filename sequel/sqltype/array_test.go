package sqltype

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Exercise the Scanner/Valuer contract, including reuse of a populated destination.
func checkRoundTrip[A ~[]E, E any, P interface {
	*A
	sql.Scanner
	driver.Valuer
}](t *testing.T, want A, encoded string) {
	t.Helper()
	var got A
	value, err := P(&got).Value()
	require.NoError(t, err)
	require.Nil(t, value)
	for _, src := range []any{encoded, []byte(encoded)} {
		require.NoError(t, P(&got).Scan(src))
		require.Equal(t, want, got)
		value, err = P(&got).Value()
		require.NoError(t, err)
		require.Equal(t, encoded, value)
	}
	require.Error(t, P(&got).Scan(42))
	require.Equal(t, want, got)
	require.NoError(t, P(&got).Scan("[]"))
	require.NotNil(t, got)
	require.Empty(t, got)
	value, err = P(&got).Value()
	require.NoError(t, err)
	require.Equal(t, "[]", value)
	require.NoError(t, P(&got).Scan(nil))
	require.Nil(t, got)
	require.NoError(t, P(&got).Scan("[]"))
	require.NotNil(t, got)
	require.Empty(t, got)
}

func TestArrayRoundTrips(t *testing.T) {
	t.Run("Int", func(t *testing.T) { checkRoundTrip(t, IntSlice[int]{-12, 0, 34}, "[-12,0,34]") })
	t.Run("Int8", func(t *testing.T) { checkRoundTrip(t, Int8Slice[int8]{-12, 0, 34}, "[-12,0,34]") })
	t.Run("Int16", func(t *testing.T) { checkRoundTrip(t, Int16Slice[int16]{-12, 0, 34}, "[-12,0,34]") })
	t.Run("Int32", func(t *testing.T) { checkRoundTrip(t, Int32Slice[int32]{-12, 0, 34}, "[-12,0,34]") })
	t.Run("Int64", func(t *testing.T) { checkRoundTrip(t, Int64Slice[int64]{-12, 0, 34}, "[-12,0,34]") })
	t.Run("Uint", func(t *testing.T) { checkRoundTrip(t, UintSlice[uint]{0, 12, 34}, "[0,12,34]") })
	t.Run("Uint8", func(t *testing.T) { checkRoundTrip(t, Uint8Slice[uint8]{0, 12, 34}, "[0,12,34]") })
	t.Run("Uint16", func(t *testing.T) { checkRoundTrip(t, Uint16Slice[uint16]{0, 12, 34}, "[0,12,34]") })
	t.Run("Uint32", func(t *testing.T) { checkRoundTrip(t, Uint32Slice[uint32]{0, 12, 34}, "[0,12,34]") })
	t.Run("Uint64", func(t *testing.T) { checkRoundTrip(t, Uint64Slice[uint64]{0, 12, 34}, "[0,12,34]") })
	t.Run("Float32", func(t *testing.T) { checkRoundTrip(t, Float32Slice[float32]{-1.5, 0, 2.25}, "[-1.5,0,2.25]") })
	t.Run("Float64", func(t *testing.T) { checkRoundTrip(t, Float64Slice[float64]{-1.5, 0, 2.25}, "[-1.5,0,2.25]") })
	t.Run("Bool", func(t *testing.T) { checkRoundTrip(t, BoolSlice[bool]{true, false, true}, "[t,f,t]") })
	t.Run("String", func(t *testing.T) {
		checkRoundTrip(t, StringSlice[string]{"", "NULL", "a,b", `a"b`, `a\b`, "世界"}, `["","NULL","a,b","a\"b","a\\b","世界"]`)
	})
}

func TestArrayScanErrors(t *testing.T) {
	factories := []func() sql.Scanner{
		func() sql.Scanner { return &IntSlice[int]{} },
		func() sql.Scanner { return &UintSlice[uint]{} },
		func() sql.Scanner { return &Float32Slice[float32]{} },
		func() sql.Scanner { return &Float64Slice[float64]{} },
		func() sql.Scanner { return &BoolSlice[bool]{} },
		func() sql.Scanner { return &StringSlice[string]{} },
	}
	for _, factory := range factories {
		t.Run(fmt.Sprintf("%T", factory()), func(t *testing.T) {
			for _, src := range []string{"", "wrong", "[", "[1", "[1,]", "[,1]", `["unterminated]`, "[]trailing", "[[1],[2]]", "[NULL]"} {
				t.Run(src, func(t *testing.T) { require.Error(t, factory().Scan(src)) })
			}
		})
	}
}

func TestInvalidElementsPreserveDestination(t *testing.T) {
	ints := Int64Slice[int64]{7}
	require.ErrorContains(t, ints.Scan("[1,nope]"), "index 1")
	require.Equal(t, Int64Slice[int64]{7}, ints)
	floats := Float64Slice[float64]{7}
	require.ErrorContains(t, floats.Scan("[1,nope]"), "index 1")
	require.Equal(t, Float64Slice[float64]{7}, floats)
	for _, input := range []string{"[t,x]", "[true]"} {
		bools := BoolSlice[bool]{true}
		require.Error(t, bools.Scan(input))
		require.Equal(t, BoolSlice[bool]{true}, bools)
	}
}

func TestUintSliceRange(t *testing.T) {
	var a Uint8Slice[uint8]
	for _, input := range []string{"[-1]", "[256]", "[abc]"} {
		require.ErrorContains(t, a.Scan(input), "index 0")
	}
	checkRoundTrip(t, Uint64Slice[uint64]{0, 18446744073709551615}, "[0,18446744073709551615]")
}

func TestSliceWhitespaceAndNamedTypes(t *testing.T) {
	type identifier int64
	var ids Int64Slice[identifier]
	require.NoError(t, ids.Scan("[ 1, -2 ,3 ]"))
	require.Equal(t, Int64Slice[identifier]{1, -2, 3}, ids)
	var floats Float32Slice[float32]
	require.NoError(t, floats.Scan("[ 1.5, -2.25 ]"))
	require.Equal(t, Float32Slice[float32]{1.5, -2.25}, floats)
	var bools BoolSlice[bool]
	require.NoError(t, bools.Scan("[ t, f ]"))
	require.Equal(t, BoolSlice[bool]{true, false}, bools)
}
