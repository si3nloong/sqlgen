package fileutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetpwd(t *testing.T) {
	require.NotEmpty(t, Getpwd())
}

func TestIsFileExists(t *testing.T) {
	require.True(t, IsFileExists("./fileutil.go"))
	require.False(t, IsFileExists("./config.abc"))
}

func TestCurDir(t *testing.T) {
	require.NotEmpty(t, CurDir())
	require.True(t, strings.HasSuffix(CurDir(), "sqlgen/internal/fileutil"))
}
