package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.True(t, *cfg.Strict)
	require.ElementsMatch(t, []string{"./**/*"}, cfg.Source)
	require.Equal(t, MySQL, cfg.Driver)
	require.Equal(t, SnakeCase, cfg.NamingConvention)
	require.Equal(t, DefaultStructTag, cfg.Tag)
	require.Equal(t, DefaultGeneratedFile, cfg.Exec.Filename)

	require.False(t, cfg.SkipHeader)
	require.False(t, cfg.SkipModTidy)
	require.False(t, cfg.SourceMap)
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
	require.Equal(t, Sqlite, cfg.Driver)
	// require.False(t, cfg.NoStrict)
}
