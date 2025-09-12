package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/cmd/codegen"
	"github.com/si3nloong/sqlgen/cmd/internal/fileutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
)

var (
	initOpts struct {
		force bool
	}
	initCmd = &cobra.Command{
		Use:   "init",
		Short: fmt.Sprintf("Set up a new %q file", codegen.DefaultConfigFile),
		RunE:  runInitCommand,
	}
)

func runInitCommand(cmd *cobra.Command, args []string) error {
	var (
		filename  = "sqlgen.yml"
		fileDest  = filepath.Join(fileutil.Getpwd(), filename)
		questions = []*survey.Question{
			{
				Name: "driver",
				Prompt: &survey.Select{
					Message: "What is your sql driver:",
					Options: []string{string(codegen.MySQL), string(codegen.Postgres), string(codegen.Sqlite)},
					Default: string(codegen.MySQL),
				},
			},
			{
				Name: "naming_convention",
				Prompt: &survey.Select{
					Message: "What is your naming convention:",
					Options: []string{string(codegen.SnakeCase), string(codegen.CamelCase), string(codegen.PascalCase)},
					Default: string(codegen.SnakeCase),
				},
			},
			{
				Name: "tag",
				Prompt: &survey.Input{
					Message: "What is your tag for parsing:",
					Default: codegen.DefaultStructTag,
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
	)

	var answer struct {
		SqlDriver        string `survey:"driver"`
		NamingConvention string `survey:"naming_convention"`
		Tag              string `survey:"tag"`
		Strict           bool   `survey:"strict"`
	}

	_, err := os.Stat(fileDest)
	if !initOpts.force && !os.IsNotExist(err) {
		cmd.Println(`Configuration file already exists`)
		return nil
	}

	if err := survey.Ask(questions, &answer); err != nil {
		return noInterruptError(err)
	}

	cfg := codegen.DefaultConfig()
	switch answer.SqlDriver {
	case string(codegen.MySQL):
		cfg.Driver = codegen.MySQL
	case string(codegen.Postgres):
		cfg.Driver = codegen.Postgres
	case string(codegen.Sqlite):
		cfg.Driver = codegen.Sqlite
	default:
		cfg.Driver = codegen.SqlDriver(answer.SqlDriver)
	}
	switch answer.NamingConvention {
	case string(codegen.SnakeCase):
		cfg.NamingConvention = codegen.SnakeCase
	case string(codegen.PascalCase):
		cfg.NamingConvention = codegen.PascalCase
	case string(codegen.CamelCase):
		cfg.NamingConvention = codegen.CamelCase
	}
	cfg.Tag = answer.Tag
	cfg.Strict = &answer.Strict

	cmd.Println("\nAbout to write to " + fileDest + ":\n")

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

	cmd.Println(`Creating ` + filename)
	return codegen.Init(cfg)
}

func init() {
	initCmd.Flags().BoolVarP(&initOpts.force, "force", "f", false, "force to execute")
}

// noInterruptError returns error when it's not `terminal.InterruptErr`
func noInterruptError(err error) error {
	if err == terminal.InterruptErr {
		return nil
	}
	return err
}
