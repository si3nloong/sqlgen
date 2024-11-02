package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesArray(t *testing.T) {

	t.Run("Scan with non-overflow & equal length string", func(t *testing.T) {
		const size = 10
		fsBytes := [size]byte{}
		bytes := FixedSizeBytes(fsBytes[:], size)
		require.Equal(t, size, bytes.size)

		rawValue := []byte(`0123456789`)
		require.NoError(t, bytes.Scan(rawValue))
		require.Equal(t, size, len(bytes.v))
		require.Equal(t, [size]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}, fsBytes)
	})

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const size = 10
		fsBytes := [size]byte{}
		require.Equal(t, fsBytes, [10]byte{})
		bytes := FixedSizeBytes(fsBytes[:], size)
		require.Equal(t, size, bytes.size)

		rawValue := []byte(`aBjk`)
		require.NoError(t, bytes.Scan(rawValue))
		require.Equal(t, size, len(bytes.v))
		require.Equal(t, [size]byte{'a', 'B', 'j', 'k'}, fsBytes)
	})

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const size = 10
		fsBytes := [size]byte{}
		require.Equal(t, fsBytes, [10]byte{})
		bytes := FixedSizeBytes(fsBytes[:], size)
		require.Equal(t, size, bytes.size)

		require.NoError(t, bytes.Scan(`Si3nloong`))
		require.Equal(t, size, len(bytes.v))
		require.Equal(t, [size]byte{'S', 'i', '3', 'n', 'l', 'o', 'o', 'n', 'g'}, fsBytes)
	})

	t.Run("Scan with overflow bytes", func(t *testing.T) {
		const size = 10
		fsBytes := [size]byte{}
		bytes := FixedSizeBytes(fsBytes[:], size)
		require.Equal(t, size, bytes.size)

		require.Error(t, bytes.Scan([]byte(`hello world`)))
	})

	t.Run("Scan with overflow string", func(t *testing.T) {
		const size = 10
		fsBytes := [size]byte{}
		bytes := FixedSizeBytes(fsBytes[:], size)
		require.Equal(t, size, bytes.size)

		require.Error(t, bytes.Scan(`hello world`))
		// require.Equal(t, fsBytes, [10]byte([]byte(`hello world`)))
	})
}
