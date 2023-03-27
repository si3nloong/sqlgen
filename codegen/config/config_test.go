package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsFileExists(t *testing.T) {
	require.True(t, isFileExists("./config.go"))
	require.False(t, isFileExists("./config.abc"))
}
