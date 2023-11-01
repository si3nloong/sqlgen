package examples_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	if err := codegen.Generate(&config.Config{
		Source: []string{"./testcase/**/*.go"},
	}); err != nil {
		panic(err)
	}

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

		// Read result file
		expected, err := os.ReadFile(path + ".gotpl")
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
