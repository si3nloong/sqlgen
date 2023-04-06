package tests_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	// Re-generate all files
	if err := filepath.Walk(".", func(path string, info fs.FileInfo, e error) error {
		if e != nil {
			return e
		}

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

		t.Run(path, func(t *testing.T) {
			require.Equal(t, expected, actual)
		})

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
