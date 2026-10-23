package encoding

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringScanner(t *testing.T) {

	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var strVal string
			str := StringScanner[string](&strVal)
			require.NoError(t, str.Scan(nil))
			require.Equal(t, "", strVal)
		})

		t.Run("string", func(t *testing.T) {
			var strVal string
			str := StringScanner[string](&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, "hello world", strVal)
		})

		t.Run("[]byte", func(t *testing.T) {
			var strVal string
			str := StringScanner[string](&strVal)
			require.NoError(t, str.Scan([]byte(`hello world`)))
			require.Equal(t, "hello world", strVal)
		})
	})

	t.Run("Scan with custom types", func(t *testing.T) {
		t.Run("custom string", func(t *testing.T) {
			type Cstring string

			var strVal Cstring
			str := StringScanner[Cstring](&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, Cstring("hello world"), strVal)
		})

		t.Run("custom []byte", func(t *testing.T) {
			type CBytes []byte

			var strVal CBytes
			str := StringScanner[CBytes](&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, CBytes("hello world"), strVal)
		})
	})
}

func TestStringScannerConversions(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	checkScanner(t, func(v *string) sql.Scanner { return StringScanner[string](v) }, func(v **string) sql.Scanner { return StringScanner[string](v) },
		[]scanCase[string]{{"hello", "hello"}, {[]byte("world"), "world"}, {true, "true"}, {int64(-42), "-42"}, {1.25, "1.25"}, {now, now.String()}}, []any{struct{}{}})
	var s string
	require.Error(t, StringScanner[string](&s, true).Scan(42.0))
	require.NoError(t, StringScanner[string](&s, true).Scan([]byte("ok")))
	require.Equal(t, "ok", s)
}
