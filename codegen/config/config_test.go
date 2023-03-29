package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, ".", cfg.SrcDir)
	require.Equal(t, "mysql", cfg.Driver)
	require.Equal(t, "sql", cfg.Tag)
	require.True(t, *cfg.IncludeHeader)
	require.Equal(t, "snake_case", cfg.NamingConvention)
}

func TestFindCfgInDir(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		f, found := findCfgInDir(".")
		require.False(t, found)
		require.Empty(t, f)
	})

	t.Run("found", func(t *testing.T) {
		f, found := findCfgInDir("../../")
		require.True(t, found)
		require.NotEmpty(t, f)
	})
}
