package examples

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/si3nloong/sqlgen/codegen/config"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	const rootDir = "./testcase"

	if err := codegen.Generate(&config.Config{
		Source:     []string{rootDir + "/**/*.go", rootDir + "/db/*"},
		SkipHeader: true,
		Exec: config.ExecConfig{
			SkipEmpty: false,
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Re-generate all files
	if err := generateModel(t, rootDir); err != nil {
		t.Fatal(err)
	}
}

func generateModel(t *testing.T, rootDir string) error {
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
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

		t.Run(fmt.Sprintf("Test %q is correct", path), func(t *testing.T) {
			require.Equal(t, expected, actual)
		})

		return nil
	})
}
