package codegen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"golang.org/x/tools/imports"
)

// Compile and exercise the generated code, not just the template parser.
func TestDatabaseTemplate(t *testing.T) {
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	write := func(path string, data []byte) {
		t.Helper()
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(dir, "go.mod"), []byte("module dbtest\n\ngo 1.27\n\nrequire github.com/si3nloong/sqlgen v0.0.0\nreplace github.com/si3nloong/sqlgen => "+strconv.Quote(root)+"\n"))
	for _, dialect := range []string{"mysql", "sqlite", "postgres"} {
		dest := filepath.Join(dir, dialect)
		if err := os.Mkdir(dest, 0700); err != nil {
			t.Fatal(err)
		}
		var importPaths []string
		funcs := template.FuncMap{
			"driver":        func() string { return dialect },
			"reserveImport": func(p string) string { importPaths = append(importPaths, strconv.Quote(p)); return "" },
			"quote":         strconv.Quote,
			"quoteVar": func(i int) string {
				if dialect == "postgres" {
					return fmt.Sprintf("$%d", i)
				}
				return "?"
			},
			"isStaticVar": func() bool { return dialect != "postgres" },
			"varRune": func() string {
				if dialect == "postgres" {
					return "$"
				}
				return "?"
			},
		}
		tmpl, err := template.New("db.go.tpl").Funcs(funcs).ParseFS(codegenTemplates, "templates/db.go.tpl")
		if err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		if err = tmpl.Execute(&b, nil); err != nil {
			t.Fatal(err)
		}
		source := []byte("package db\nimport (\n" + strings.Join(importPaths, "\n") + "\n)\n" + b.String())
		source, err = imports.Process(filepath.Join(dest, "db.go"), source, nil)
		if err != nil {
			t.Fatalf("%s: %v", dialect, err)
		}
		write(filepath.Join(dest, "db.go"), source)
		tests, err := template.New("db_regression_test.go.tpl").Delims("[[", "]]").Funcs(funcs).ParseFiles("testdata/db_regression_test.go.tpl")
		if err != nil {
			t.Fatal(err)
		}
		b.Reset()
		if err = tests.Execute(&b, nil); err != nil {
			t.Fatal(err)
		}
		write(filepath.Join(dest, "db_test.go"), b.Bytes())
	}
	cmd := exec.Command("go", "test", "-mod=mod", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated database regression tests: %v\n%s", err, output)
	}
}
