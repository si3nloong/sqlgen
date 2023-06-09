package cmd

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/internal/fileutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Set up a new or existing `sqlgen.yaml` file.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(`This utility will walk you through creating a sqlgen.yaml file.
It only covers the most common items, and tries to guess sensible defaults.

See ` + "`sqlgen init`" + ` for definitive documentation on these fields
and exactly what they do.`)
			return nil
		},
		RunE: runInitCommand,
	}
)

func runInitCommand(cmd *cobra.Command, args []string) error {
	var (
		questions = []*survey.Question{
			{
				Name: "driver",
				Prompt: &survey.Select{
					Message: "What is your sql driver:",
					Options: []string{"mysql", "postgres", "sqlite", "sql"},
					Default: "mysql",
				},
			},
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
					Message: "What is your tag for parsing:",
					Default: "sql",
				},
			},
			{
				Name: "strict",
				Prompt: &survey.Confirm{
					Message: "Is it strict parsing:",
					Default: true,
				},
			},
		}
		answers = config.DefaultConfig()
	)

	if err := survey.Ask(questions, answers); err != nil {
		return noInterruptError(err)
	}

	w := bytes.NewBufferString("")
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	if err := enc.Encode(answers); err != nil {
		return err
	}

	fileDest := filepath.Join(fileutil.Getpwd(), "sqlgen.yml")
	cmd.Println("\nAbout to write to " + fileDest + ":\n")
	cmd.Println(w.String())

	var ok bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Is this OK?",
		Default: true,
	}, &ok); err != nil {
		return noInterruptError(err)
	}

	// Do nothing when user choose "No".
	if !ok {
		return nil
	}

	if err := os.WriteFile(fileDest, w.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

// noInterruptError returns error when it's not `terminal.InterruptErr`
func noInterruptError(err error) error {
	if err == terminal.InterruptErr {
		return nil
	}
	return err
}
