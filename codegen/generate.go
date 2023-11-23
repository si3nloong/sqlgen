package codegen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/tools/imports"
)

type Generator struct {
	dialect   sequel.Dialect
	quoteChar rune
}

func (g Generator) QuoteStart() string {
	return string(g.quoteChar)
}

func (g Generator) Quote(v string) string {
	return string(g.quoteChar) + v + string(g.quoteChar)
}

func (g Generator) QuoteEnd() string {
	return string(g.quoteChar)
}

func Init(cfg *config.Config) error {
	tmpl, err := template.ParseFS(codegenTemplates, "templates/init.yml.go.tpl")
	if err != nil {
		return err
	}

	w, err := os.OpenFile(config.DefaultConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := tmpl.Execute(w, cfg); err != nil {
		return err
	}
	return nil
}

func renderTemplate[T templates.ModelTmplParams | struct{}](
	tmplName string,
	skipHeader bool,
	dialect sequel.Dialect,
	pkgPath string,
	pkgName string,
	getter string,
	dstDir string,
	dstFilename string,
	params T,
) error {
	w, blr := strpool.AcquireString(), strpool.AcquireString()
	defer func() {
		strpool.ReleaseString(w)
		strpool.ReleaseString(blr)
	}()

	quoteChar := rune('"')
	switch dialect.QuoteChar() {
	case '`':
		quoteChar = '"'
	case '"':
		quoteChar = '`'
	}

	g := &Generator{quoteChar: quoteChar}
	impPkg := NewPackage(pkgPath, pkgName)
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{
		"quote":             g.Quote,
		"createTable":       g.createTableStmt(dialect),
		"alterTable":        alterTableStmt(dialect),
		"insertOneStmt":     g.insertOneStmt(dialect),
		"findByPKStmt":      g.findByPKStmt(dialect),
		"updateByPKStmt":    g.updateByPKStmt(dialect),
		"reserveImport":     reserveImport(impPkg),
		"castAs":            castAs(impPkg),
		"addrOf":            addrOf(impPkg),
		"wrap":              dialect.Wrap,
		"getFieldTypeValue": getFieldTypeValue(impPkg, getter),
		"varStmt":           varStmt(dialect),
		"var":               dialect.Var,
		"dialectVar":        dialectVar(dialect),
	}).ParseFS(codegenTemplates, "templates/"+tmplName)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(blr, params); err != nil {
		return err
	}

	if !skipHeader {
		w.WriteString(fmt.Sprintf("// Code generated by sqlgen, version %s. DO NOT EDIT.\n\n", sequel.Version))
	}

	w.WriteString("package " + pkgName + "\n\n")

	if len(impPkg.imports) > 0 {
		w.WriteString("import (\n")
		for _, pkg := range impPkg.imports {
			if filepath.Base(pkg.Path()) == pkg.Name() {
				w.WriteString("\t" + strconv.Quote(pkg.Path()) + "\n")
			} else {
				w.WriteString("\t" + pkg.Name() + " " + strconv.Quote(pkg.Path()) + "\n")
			}
		}
		w.WriteString(")\n")
	}
	w.WriteString(blr.String())

	os.MkdirAll(dstDir, fileMode)
	fileDest := filepath.Join(dstDir, dstFilename)
	formatted, err := imports.Process(fileDest, []byte(w.String()), &imports.Options{Comments: true})
	if err != nil {
		return err
	}
	blr.Reset()
	w.Reset()

	log.Println("Creating " + fileDest)
	if err := os.WriteFile(fileDest, formatted, fileMode); err != nil {
		return err
	}
	return nil
}
