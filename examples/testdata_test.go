package examples_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/si3nloong/sqlgen/codegen/config"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/sqlite"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	if err := codegen.Generate(&config.Config{
		Source: []string{"./testcase/**/*.go"},
	}); err != nil {
		panic(err)
	}

	// Re-generate all files
	if err := filepath.Walk("./testcase", func(path string, info fs.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !info.IsDir() {
			return nil
		}

		if fileutil.IsDirEmptyFiles(path) {
			return nil
		}

		actual, err := os.ReadFile(filepath.Join(path, config.DefaultGeneratedFile))
		if err != nil {
			return err
		}

		// Read result file
		expected, err := os.ReadFile(filepath.Join(path, config.DefaultGeneratedFile+".tpl"))
		if err != nil {
			return fmt.Errorf("%w, happened in directory %q", err, path)
		}

		t.Run("Compare the []byte in directory "+path, func(t *testing.T) {
			require.Equal(t, expected, actual)
		})

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
