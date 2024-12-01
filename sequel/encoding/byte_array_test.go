package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesArray(t *testing.T) {

	t.Run("Scan with non-overflow & equal length string", func(t *testing.T) {
		const SIZE = 10
		fsBytes := [SIZE]byte{'0', 'V'}
		bytes := ByteArrayScanner(fsBytes[:], SIZE)
		require.NoError(t, bytes.Scan(nil))
		require.Equal(t, [SIZE]byte{}, fsBytes)

		rawValue := []byte(`0123456789`)
		require.NoError(t, bytes.Scan(rawValue))
		require.Equal(t, [SIZE]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}, fsBytes)
	})

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const SIZE = 10
		fsBytes := [SIZE]byte{}
		require.Equal(t, fsBytes, [10]byte{})
		bytes := ByteArrayScanner(fsBytes[:], SIZE)

		rawValue := []byte(`aBjk`)
		require.NoError(t, bytes.Scan(rawValue))
		require.Equal(t, [SIZE]byte{'a', 'B', 'j', 'k'}, fsBytes)
	})

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const SIZE = 10
		fsBytes := [SIZE]byte{}
		require.Equal(t, fsBytes, [10]byte{})
		bytes := ByteArrayScanner(fsBytes[:], SIZE)

		require.NoError(t, bytes.Scan(`Si3nloong`))
		require.Equal(t, [SIZE]byte{'S', 'i', '3', 'n', 'l', 'o', 'o', 'n', 'g'}, fsBytes)
	})

	t.Run("Scan with overflow bytes", func(t *testing.T) {
		const SIZE = 10
		fsBytes := [SIZE]byte{}
		bytes := ByteArrayScanner(fsBytes[:], SIZE)

		require.Error(t, bytes.Scan([]byte(`hello world`)))
	})

	t.Run("Scan with overflow string", func(t *testing.T) {
		const SIZE = 10
		fsBytes := [SIZE]byte{}
		bytes := ByteArrayScanner(fsBytes[:], SIZE)

		require.Error(t, bytes.Scan(`hello world`))
	})
}
