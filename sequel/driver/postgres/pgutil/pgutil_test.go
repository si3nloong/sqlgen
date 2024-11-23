package pgutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuote(t *testing.T) {
	t.Run("With chinese", func(t *testing.T) {
		require.Equal(t, Quote(`Testing哈咯`), `E'Testing哈咯'`)
	})

	t.Run("With quote", func(t *testing.T) {
		require.Equal(t, Quote(`Te'st'ing`), `E'Te''st''ing'`)
		require.Equal(t, Quote(`Test'ing`), `E'Test''ing'`)
	})

	t.Run("With breakline", func(t *testing.T) {
		require.Equal(t, Quote(`Test'
	
	
		ing`), `E'Test''\n\t\n\t\n\t\ting'`)
	})

	t.Run("With backslash", func(t *testing.T) {
		require.Equal(t, Quote(`	Test'
\ing`), `E'\tTest''\n\\ing'`)
	})
}
