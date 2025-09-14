package encoding

import (
	"testing"

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
