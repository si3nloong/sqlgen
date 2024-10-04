package enum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum(t *testing.T) {
	values := (Custom{}).Values()
	require.ElementsMatch(t, []any{"", int64(0), int64(0)}, values)
	// require.ElementsMatch(t, []any{"", int64(1), int64(0)}, values)
}
