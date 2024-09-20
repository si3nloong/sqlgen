package examples

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	const rootDir = "./testcase"

	if err := codegen.Generate(&codegen.Config{
		Source:     []string{rootDir + "/**/*.go"},
		SkipHeader: true,
		// Driver:     codegen.Postgres,
		Database: &codegen.DatabaseConfig{
			Package: "mysqldb",
			Dir:     "./db/mysql",
		},
		// QuoteIdentifier: true,
		Exec: codegen.ExecConfig{
			SkipEmpty: false,
		},
		DataTypes: map[string]codegen.DataType{
			"github.com/paulmach/orb.Point": {
				DataType:   "POINT NOT NULL",
				SQLScanner: `ST_AsBinary({{.}},4326)`,
				Scanner:    `github.com/paulmach/orb/encoding/ewkb.Scanner({{.}})`,
				SQLValuer:  `ST_GeomFromEWKB({{.}})`,
				Valuer:     `github.com/paulmach/orb/encoding/ewkb.Value({{.}},4326)`,
			},
			"github.com/gofrs/uuid/v5.UUID": {
				DataType: "UUID",
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

	if err := codegen.Generate(&codegen.Config{
		Source:     []string{},
		SkipHeader: true,
		Driver:     codegen.Postgres,
		Database: &codegen.DatabaseConfig{
			Package: "postgresdb",
			Dir:     "./db/postgres",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := codegen.Generate(&codegen.Config{
		Source:     []string{},
		SkipHeader: true,
		Driver:     codegen.Sqlite,
		Database: &codegen.DatabaseConfig{
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

		actual, err := os.ReadFile(filepath.Join(path, codegen.DefaultGeneratedFile))
		if err != nil {
			return err
		}

		// Read result file
		expected, err := os.ReadFile(filepath.Join(path, codegen.DefaultGeneratedFile+".tpl"))
		if err != nil {
			return fmt.Errorf("%w, happened in directory %q", err, path)
		}

		t.Run(fmt.Sprintf("Test %q is correct", path), func(t *testing.T) {
			require.Equal(t, expected, actual)
		})

		return nil
	})
}
