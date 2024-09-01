package codegen

import (
	"bytes"
	"go/types"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/imports"
)

func Init(cfg *Config) error {
	f, err := os.OpenFile(DefaultConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewEncoder(f, yaml.WithComment(yaml.CommentMap{
		"$.src":               []*yaml.Comment{yaml.HeadComment(" Where are all the model files located? Globs are supported eg: src/**/*.go")},
		"$.driver":            []*yaml.Comment{yaml.HeadComment(` Optional: possibly values : "mysql", "postgres" or "sqlite". Default value is "mysql".`)},
		"$.naming_convention": []*yaml.Comment{yaml.HeadComment(` Optional: possibly values : "snake_case", "camelCase" or "PascalCase". Default value is "snake_case".`)},
		"$.struct_tag":        []*yaml.Comment{yaml.HeadComment(` Optional: the struct tag for "sqlgen" to read from. Default is "sql".`)},
		"$.quote_identifier":  []*yaml.Comment{yaml.HeadComment(` Optional: whether to quote the table or column name. Default value is "false"`)},
		"$.strict":            []*yaml.Comment{yaml.HeadComment(` Optional: enables a wide range of type checking behavior that results in stronger guarantees of program correctness. Default value is "true".`)},
		"$.exec":              []*yaml.Comment{yaml.HeadComment(` Optional: where should the generated model code go?`)},
		"$.exec.skip_empty":   []*yaml.Comment{yaml.HeadComment(` Optional: whether to not generate codes for struct that has no field.`)},
		"$.database":          []*yaml.Comment{yaml.HeadComment(` Optional: where should the generated database code go?`)},
		"$.source_map":        []*yaml.Comment{yaml.HeadComment(` Optional: generate source map. Default value is "false".`)},
		"$.skip_mod_tidy":     []*yaml.Comment{yaml.HeadComment(` Optional: set to skip running "go mod tidy" when generating server code.`)},
		"$.skip_header":       []*yaml.Comment{yaml.HeadComment(` Optional: turn on to not generate any file header in generated files.`)},
		"$.data_types":        []*yaml.Comment{yaml.HeadComment(` Optional: configure column type mapping for go types.`)},
	})).Encode(cfg); err != nil {
		return err
	}

	return f.Close()
}

func renderTemplate(
	g *Generator,
	tmplName string,
	pkgPath string,
	pkgName string,
	dstDir string,
	dstFilename string,
) error {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)

	impPkg := NewPackage(pkgPath, pkgName)
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{
		"driver":   g.dialect.Driver,
		"quote":    g.Quote,
		"quoteVar": g.dialect.QuoteVar,
		"isStaticVar": func() bool {
			return g.dialect.QuoteVar(1) == g.dialect.QuoteVar(2)
		},
		"reserveImport": reserveImport(impPkg),
		"varRune": func() string {
			return string(g.dialect.VarRune())
		},
	}).ParseFS(codegenTemplates, "templates/"+tmplName)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(blr, struct{}{}); err != nil {
		return err
	}

	g.Buffer = new(bytes.Buffer)
	g.buildHeader()
	g.L("package " + pkgName)
	g.L()

	if len(impPkg.imports) > 0 {
		g.L("import (")
		for _, pkg := range impPkg.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				g.L(strconv.Quote(pkg.Path()))
			} else {
				g.L(pkg.Name() + " " + strconv.Quote(pkg.Path()))
			}
		}
		g.L(")")
	}

	g.WriteString(blr.String())
	strpool.ReleaseString(blr)

	os.MkdirAll(dstDir, fileMode)
	fileDest := filepath.Join(dstDir, dstFilename)
	// formatted, err := format.Source([]byte(w.String()))
	// if err != nil {
	// 	return err
	// }
	formatted, err := imports.Process(fileDest, g.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		return err
	}
	g.Reset()

	slog.Info("Creating " + fileDest)
	if err := os.WriteFile(fileDest, formatted, fileMode); err != nil {
		return err
	}
	return nil
}

func reserveImport(impPkgs *Package) func(pkgPath string, aliases ...string) string {
	return func(pkgPath string, aliases ...string) string {
		name := filepath.Base(pkgPath)
		if len(aliases) > 0 {
			name = aliases[0]
		}
		impPkgs.Import(types.NewPackage(pkgPath, name))
		return ""
	}
}
