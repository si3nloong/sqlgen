package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunesArray(t *testing.T) {

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const size = 10
		rBytes := [size]rune{}
		runes := RuneArrayScanner(rBytes[:], size)

		rawValue := []byte(`开心🥳`)
		require.Equal(t, 10, len(rawValue))
		require.NoError(t, runes.Scan(rawValue))
		require.Equal(t, [10]rune{'开', '心', '🥳'}, rBytes)
	})

	t.Run("Scan with non-overflow & equal length string", func(t *testing.T) {
		const size = 10
		rBytes := [size]rune{}
		runes := RuneArrayScanner(rBytes[:], size)

		rawValue := []byte(`我知道现在很难成功！`)
		require.Equal(t, 30, len(rawValue))
		require.NoError(t, runes.Scan(rawValue))
		require.Equal(t, [10]rune{'我', '知', '道', '现', '在', '很', '难', '成', '功', '！'}, rBytes)
	})

	t.Run("Scan with overflow string", func(t *testing.T) {
		const size = 5
		rBytes := [size]rune{}
		runes := RuneArrayScanner(rBytes[:], size)

		rawValue := []byte(`现在最夯的产品！🔥`)
		require.Equal(t, 28, len(rawValue))
		require.Error(t, runes.Scan(rawValue))
	})
}

func TestRuneArrayStringsAndReset(t *testing.T) {
	for _, input := range []string{"abc", "开心🥳"} {
		t.Run(input, func(t *testing.T) {
			var dst [3]rune
			scanner := RuneArrayScanner(dst[:], len(dst))
			require.NoError(t, scanner.Scan(input))
			require.Equal(t, []rune(input), dst[:])
			require.Error(t, scanner.Scan(input+"x"))
			require.Error(t, scanner.Scan(42))
			require.NoError(t, scanner.Scan(nil))
			require.Equal(t, [3]rune{}, dst)
			require.NoError(t, scanner.Scan("xyz"))
			require.Equal(t, [3]rune{'x', 'y', 'z'}, dst)
		})
	}
}
