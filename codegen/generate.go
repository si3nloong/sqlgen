package codegen

import (
	"os"
	"text/template"

	"github.com/si3nloong/sqlgen/codegen/config"
)

func Init(cfg *config.Config) error {
	tmpl, err := template.ParseFS(codegenTemplates, "templates/init.yml.go.tpl")
	if err != nil {
		return err
	}

	w, err := os.OpenFile(config.DefaultConfigFile, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := tmpl.Execute(w, cfg); err != nil {
		return err
	}
	return nil
}
