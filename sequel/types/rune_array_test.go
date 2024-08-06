package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunesArray(t *testing.T) {

	t.Run("Scan with non-overflow & smaller length string", func(t *testing.T) {
		const size = 10
		rBytes := [size]rune{}
		runes := FixedSizeRunes(rBytes[:], size)
		require.Equal(t, size, runes.size)

		rawValue := []byte(`开心🥳`)
		require.Equal(t, 10, len(rawValue))
		require.NoError(t, runes.Scan(rawValue))
		require.Equal(t, size, len(runes.v))
		require.Equal(t, [10]rune{'开', '心', '🥳'}, rBytes)
	})

	t.Run("Scan with non-overflow & equal length string", func(t *testing.T) {
		const size = 10
		rBytes := [size]rune{}
		runes := FixedSizeRunes(rBytes[:], size)
		require.Equal(t, size, runes.size)

		rawValue := []byte(`我知道现在很难成功！`)
		require.Equal(t, 30, len(rawValue))
		require.NoError(t, runes.Scan(rawValue))
		require.Equal(t, size, len(runes.v))
		require.Equal(t, [10]rune{'我', '知', '道', '现', '在', '很', '难', '成', '功', '！'}, rBytes)
	})

	t.Run("Scan with overflow string", func(t *testing.T) {
		const size = 5
		rBytes := [size]rune{}
		runes := FixedSizeRunes(rBytes[:], size)
		require.Equal(t, size, runes.size)

		rawValue := []byte(`现在最夯的产品！🔥`)
		require.Equal(t, 28, len(rawValue))
		require.Error(t, runes.Scan(rawValue))
	})
}
