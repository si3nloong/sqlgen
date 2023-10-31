package alias

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAliasStruct(t *testing.T) {
	require.Equal(t, float64(0), (&AliasStruct{}).a)
}
