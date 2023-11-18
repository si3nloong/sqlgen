package examples

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/sequel"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/mysql"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/postgres"
	_ "github.com/si3nloong/sqlgen/sequel/dialect/sqlite"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/internal/fileutil"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	const rootDir = "./testcase"

	if err := codegen.Generate(&config.Config{
		Source:     []string{rootDir + "/**/*.go"},
		SkipHeader: true,
	}); err != nil {
		t.Fatal(err)
	}

	// if err := patchVersionInFiles(t, rootDir); err != nil {
	// 	t.Fatal(err)
	// }

	// Re-generate all files
	if err := generateModel(t, rootDir); err != nil {
		t.Fatal(err)
	}
}

func patchVersionInFiles(t *testing.T, rootDir string) error {
	headerRegex := regexp.MustCompile(`\/\/ Code generated by sqlgen, version (.*)\. DO NOT EDIT\.`)
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		// Read result file
		filename := filepath.Join(path, config.DefaultGeneratedFile+".tpl")
		f, err := os.OpenFile(filename, os.O_RDWR, 0o755)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return fmt.Errorf("%w, happened in directory %q", err, path)
		}
		defer f.Close()

		stat, _ := f.Stat()
		b := make([]byte, stat.Size())
		_, err = bufio.NewReader(f).Read(b)
		if err != nil {
			return err
		}

		idx := bytes.IndexByte(b, '\n')
		line := b[:idx]
		matches := headerRegex.FindSubmatch(line)
		if len(matches) > 1 {
			newVersion := []byte(sequel.Version)
			if bytes.Equal(matches[1], newVersion) {
				return nil
			}

			b = b[idx:]

			if err := f.Truncate(0); err != nil {
				return err
			}
			if _, err := f.Seek(0, 0); err != nil {
				return err
			}

			f.Write(bytes.ReplaceAll(line, matches[1], newVersion))
			f.Write(b)

			log.Println("Patching " + filename)
		}

		return nil
	})
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
