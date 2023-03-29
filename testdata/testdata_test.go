package testdata_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	if err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || filepath.Base(path) != "generated.go" {
			return nil
		}

		actual, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		expected, err := os.ReadFile(path + ".txt")
		if err != nil {
			return err
		}

		require.Equal(t, expected, actual)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
