package cmd

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/internal/tools"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "init",
		RunE: runInitCommand,
	}
}

func runInitCommand(cmd *cobra.Command, args []string) error {
	var (
		questions = []*survey.Question{
			{
				Name: "namingConvention",
				Prompt: &survey.Select{
					Message: "What is your naming convention:",
					Options: []string{"snake_case", "camelCase", "PascalCase"},
					Default: "snake_case",
				},
			},
			{
				Name: "tag",
				Prompt: &survey.Input{
					Message: "Your required tag for parsing",
					Default: "sql",
				},
			},
		}
		answers config.Config
	)

	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	w := bytes.NewBufferString(``)
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	if err := enc.Encode(answers); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(tools.Getpwd(), "sqlgen.yml"), w.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}
