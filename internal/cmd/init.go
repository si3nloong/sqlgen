package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/si3nloong/sqlgen/codegen"
	"github.com/si3nloong/sqlgen/codegen/config"
	"github.com/si3nloong/sqlgen/internal/fileutil"

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
		Short: fmt.Sprintf("Set up a new %q file", config.DefaultConfigFile),
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
					Options: []string{string(config.MySQL), string(config.Postgres), string(config.Sqlite), string(config.Clickhouse)},
					Default: string(config.MySQL),
				},
			},
			{
				Name: "naming_convention",
				Prompt: &survey.Select{
					Message: "What is your naming convention:",
					Options: []string{string(config.SnakeCase), string(config.CamelCase), string(config.PascalCase)},
					Default: string(config.SnakeCase),
				},
			},
			{
				Name: "tag",
				Prompt: &survey.Input{
					Message: "What is your tag for parsing:",
					Default: config.DefaultStructTag,
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
		Strict           bool   `survey:"strict,omitempty"`
	}

	_, err := os.Stat(fileDest)
	if !initOpts.force && !os.IsNotExist(err) {
		cmd.Println(`Configuration file already exists`)
		return nil
	}

	if err := survey.Ask(questions, &answer); err != nil {
		return noInterruptError(err)
	}

	cfg := config.DefaultConfig()
	switch answer.SqlDriver {
	case string(config.MySQL):
		cfg.Driver = config.MySQL
	case string(config.Postgres):
		cfg.Driver = config.Postgres
	case string(config.Sqlite):
		cfg.Driver = config.Sqlite
	default:
		cfg.Driver = config.SqlDriver(answer.SqlDriver)
	}
	switch answer.NamingConvention {
	case string(config.SnakeCase):
		cfg.NamingConvention = config.SnakeCase
	case string(config.PascalCase):
		cfg.NamingConvention = config.PascalCase
	case string(config.CamelCase):
		cfg.NamingConvention = config.CamelCase
	}
	cfg.Tag = answer.Tag
	cfg.Strict = answer.Strict

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
