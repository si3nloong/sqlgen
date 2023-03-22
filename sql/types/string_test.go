package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Run("Scan with primitive types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			var strVal string
			str := String(&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, "hello world", str.Interface())
		})

		t.Run("[]byte", func(t *testing.T) {
			var strVal string
			str := String(&strVal)
			require.NoError(t, str.Scan([]byte(`hello world`)))
			require.Equal(t, "hello world", str.Interface())
		})
	})

	t.Run("Scan with custom types", func(t *testing.T) {
		t.Run("custom string", func(t *testing.T) {
			type Cstring string

			var strVal Cstring
			str := String(&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, Cstring("hello world"), str.Interface())
		})

		t.Run("custom []byte", func(t *testing.T) {
			type CBytes []byte

			var strVal CBytes
			str := String(&strVal)
			require.NoError(t, str.Scan("hello world"))
			require.Equal(t, CBytes("hello world"), str.Interface())
		})
	})
}
