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
		Source:     []string{rootDir + "/**/*.go"},
		SkipHeader: true,
		Database: &config.DatabaseConfig{
			Package: "mysqldb",
			Dir:     "./db/mysql",
		},
		Exec: config.ExecConfig{
			SkipEmpty: false,
		},
		Models: map[string]*config.Model{
			"github.com/paulmach/orb.Point": {
				DataType:   "POINT",
				SQLScanner: `ST_AsBinary({{.}}, 4326)`,
				Scanner:    `github.com/paulmach/orb/encoding/ewkb.Scanner({{.}})`,
				SQLValuer:  `ST_GeomFromEWKB({{.}})`,
				Valuer:     `github.com/paulmach/orb/encoding/ewkb.Value({{.}}, 4326)`,
			},
			"encoding/json.Number": {
				DataType: "VARCHAR(20)",
				Scanner:  "github.com/si3nloong/sqlgen/examples/testcase/struct-field/json.Number({{.}})",
				Valuer:   "{{.}}.String()",
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := codegen.Generate(&config.Config{
		Source:     []string{},
		SkipHeader: true,
		Driver:     config.Postgres,
		Database: &config.DatabaseConfig{
			Package: "postgresdb",
			Dir:     "./db/postgres",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := codegen.Generate(&config.Config{
		Source:     []string{},
		SkipHeader: true,
		Driver:     config.Sqlite,
		Database: &config.DatabaseConfig{
			Package: "sqlite",
			Dir:     "./db/sqlite",
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
