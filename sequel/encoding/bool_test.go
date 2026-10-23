package encoding

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	t.Run("Scan with nil", func(t *testing.T) {
		var flag *bool
		require.Nil(t, flag)
		v := BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(true))
		require.NotNil(t, flag)
		require.True(t, *flag)

		v = BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(false))
		require.NotNil(t, flag)
		require.False(t, *flag)

		var ptrflag = new(bool)
		require.NotNil(t, ptrflag)
		v = BoolScanner[bool](&ptrflag)
		require.NoError(t, v.Scan(nil))
		require.Nil(t, ptrflag)
	})

	t.Run("Scan with bool", func(t *testing.T) {
		var flag bool
		v := BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(true))
		require.True(t, flag)

		v = BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(false))
		require.False(t, flag)
	})

	t.Run("Scan with []byte", func(t *testing.T) {
		var flag bool
		v := BoolScanner[bool](&flag)
		require.NoError(t, v.Scan([]byte(`true`)))
		require.True(t, flag)

		v = BoolScanner[bool](&flag)
		require.NoError(t, v.Scan([]byte(`false`)))
		require.False(t, flag)
	})

	t.Run("Scan with int64", func(t *testing.T) {
		var flag bool
		v := BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(int64(1)))
		require.True(t, flag)

		v = BoolScanner[bool](&flag)
		require.NoError(t, v.Scan(int64(0)))
		require.False(t, flag)
	})

	t.Run("Scan with custom bool type", func(t *testing.T) {
		type CBool bool
		var flag CBool
		v := BoolScanner[CBool](&flag)
		require.NoError(t, v.Scan(true))
		require.Equal(t, CBool(true), flag)

		v = BoolScanner[CBool](&flag)
		require.NoError(t, v.Scan(false))
		require.Equal(t, CBool(false), flag)
	})
}

func TestBoolScannerConversions(t *testing.T) {
	checkScanner(t, func(v *bool) sql.Scanner { return BoolScanner[bool](v) }, func(v **bool) sql.Scanner { return BoolScanner[bool](v) },
		[]scanCase[bool]{{true, true}, {false, false}, {int64(2), true}, {int64(0), false}, {"true", true}, {[]byte("false"), false}}, []any{"invalid", []byte("invalid"), 1.5})
	var b bool
	require.Error(t, BoolScanner[bool](&b, true).Scan("true"))
	require.NoError(t, BoolScanner[bool](&b, true).Scan(true))
	require.True(t, b)
}
