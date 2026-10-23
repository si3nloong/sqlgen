package pgtype

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
	require.NoError(t, P(&got).Scan("{}"))
	require.NotNil(t, got)
	require.Empty(t, got)
	value, err = P(&got).Value()
	require.NoError(t, err)
	require.Equal(t, "{}", value)
	require.NoError(t, P(&got).Scan(nil))
	require.Nil(t, got)
	require.NoError(t, P(&got).Scan("{}"))
	require.NotNil(t, got)
	require.Empty(t, got)
}

func TestArrayRoundTrips(t *testing.T) {
	t.Run("Int", func(t *testing.T) { checkRoundTrip(t, IntArray[int]{-12, 0, 34}, "{-12,0,34}") })
	t.Run("Int8", func(t *testing.T) { checkRoundTrip(t, Int8Array[int8]{-12, 0, 34}, "{-12,0,34}") })
	t.Run("Int16", func(t *testing.T) { checkRoundTrip(t, Int16Array[int16]{-12, 0, 34}, "{-12,0,34}") })
	t.Run("Int32", func(t *testing.T) { checkRoundTrip(t, Int32Array[int32]{-12, 0, 34}, "{-12,0,34}") })
	t.Run("Int64", func(t *testing.T) { checkRoundTrip(t, Int64Array[int64]{-12, 0, 34}, "{-12,0,34}") })
	t.Run("Uint", func(t *testing.T) { checkRoundTrip(t, UintArray[uint]{0, 12, 34}, "{0,12,34}") })
	t.Run("Uint8", func(t *testing.T) { checkRoundTrip(t, Uint8Array[uint8]{0, 12, 34}, "{0,12,34}") })
	t.Run("Uint16", func(t *testing.T) { checkRoundTrip(t, Uint16Array[uint16]{0, 12, 34}, "{0,12,34}") })
	t.Run("Uint32", func(t *testing.T) { checkRoundTrip(t, Uint32Array[uint32]{0, 12, 34}, "{0,12,34}") })
	t.Run("Uint64", func(t *testing.T) { checkRoundTrip(t, Uint64Array[uint64]{0, 12, 34}, "{0,12,34}") })
	t.Run("Float32", func(t *testing.T) { checkRoundTrip(t, Float32Array[float32]{-1.5, 0, 2.25}, "{-1.5,0,2.25}") })
	t.Run("Float64", func(t *testing.T) { checkRoundTrip(t, Float64Array[float64]{-1.5, 0, 2.25}, "{-1.5,0,2.25}") })
	t.Run("Bool", func(t *testing.T) { checkRoundTrip(t, BoolArray[bool]{true, false, true}, "{t,f,t}") })
	t.Run("String", func(t *testing.T) {
		checkRoundTrip(t, StringArray[string]{"", "NULL", "a,b", `a"b`, `a\b`, "世界"}, `{"","NULL","a,b","a\"b","a\\b","世界"}`)
	})
}

func TestArrayScanErrors(t *testing.T) {
	factories := []func() sql.Scanner{
		func() sql.Scanner { return &IntArray[int]{} },
		func() sql.Scanner { return &UintArray[uint]{} },
		func() sql.Scanner { return &Float32Array[float32]{} },
		func() sql.Scanner { return &Float64Array[float64]{} },
		func() sql.Scanner { return &BoolArray[bool]{} },
		func() sql.Scanner { return &StringArray[string]{} },
	}
	for _, factory := range factories {
		t.Run(fmt.Sprintf("%T", factory()), func(t *testing.T) {
			for _, src := range []string{"", "wrong", "{", "{1", "{1,}", "{,1}", `{"unterminated}`, "{}trailing", "{{1},{2}}", "{NULL}"} {
				t.Run(src, func(t *testing.T) { require.Error(t, factory().Scan(src)) })
			}
		})
	}
}

func TestInvalidElementsPreserveDestination(t *testing.T) {
	ints := Int64Array[int64]{7}
	require.ErrorContains(t, ints.Scan("{1,nope}"), "index 1")
	require.Equal(t, Int64Array[int64]{7}, ints)
	floats := Float64Array[float64]{7}
	require.ErrorContains(t, floats.Scan("{1,nope}"), "index 1")
	require.Equal(t, Float64Array[float64]{7}, floats)
	for _, input := range []string{"{t,x}", "{true}"} {
		bools := BoolArray[bool]{true}
		require.Error(t, bools.Scan(input))
		require.Equal(t, BoolArray[bool]{true}, bools)
	}
}

func TestByteaArray(t *testing.T) {
	checkRoundTrip(t, ByteaArray[[]byte]{{0, 255}, {}, []byte("hello")}, `{"\\x00ff","\\x","\\x68656c6c6f"}`)
	for _, src := range []string{`{"\\xzz"}`, `{"\\x0"}`, `{"\\1"}`, `{"\\999"}`, "invalid"} {
		var a ByteaArray[[]byte]
		require.Error(t, a.Scan(src))
	}
}

func TestParseBytea(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  []byte
	}{
		{`\x00ff`, []byte{0, 255}},
		{`plain`, []byte("plain")},
		{`a\000b`, []byte{'a', 0, 'b'}},
		{`a\\b`, []byte(`a\b`)},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := parseBytea([]byte(tc.input))
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
