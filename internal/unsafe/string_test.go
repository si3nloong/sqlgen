package unsafe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsafeString(t *testing.T) {
	type C string

	require.Equal(t, `xdd`, String("xdd"))
	require.Equal(t, `abcd`, String(C("abcd")))
	require.Equal(t, `hello world!`, String([]byte(`hello world!`)))
}
