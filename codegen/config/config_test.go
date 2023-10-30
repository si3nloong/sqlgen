package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.ElementsMatch(t, []string{"."}, cfg.Source)
	require.Equal(t, MySQL, cfg.Driver)
	require.Equal(t, "sql", cfg.Tag)
	require.Equal(t, "snake_case", cfg.NamingConvention)

	require.True(t, cfg.IncludeHeader)
	require.True(t, cfg.Strict)
}

func TestFindCfgInDir(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		f, found := findCfgInDir(".")
		require.False(t, found)
		require.Empty(t, f)
	})

	t.Run("found", func(t *testing.T) {
		f, found := findCfgInDir("./testdata/")
		require.True(t, found)
		require.NotEmpty(t, f)
	})
}

func TestLoadConfigFrom(t *testing.T) {
	var (
		cfg *Config
		err error
	)
	cfg, err = LoadConfigFrom(".")
	require.Error(t, err)
	require.Nil(t, cfg)

	cfg, err = LoadConfigFrom("./testdata/config.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.True(t, cfg.Strict)
}
