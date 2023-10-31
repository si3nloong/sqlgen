package cmd

import (
	"bytes"
	"log"
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
		Short: "Set up a new or existing `sqlgen.yml` file.",
		// 		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 			cmd.Println(`This utility will walk you through creating a sqlgen.yaml file.
		// It only covers the most common items, and tries to guess sensible defaults.

		// See ` + "`sqlgen init`" + ` for definitive documentation on these fields
		// and exactly what they do.`)
		// 			return nil
		// 		},
		RunE: runInitCommand,
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
					Options: []string{string(config.MySQL), string(config.Postgres), string(config.Sqlite)},
					Default: string(config.MySQL),
				},
			},
			{
				Name: "namingConvention",
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
	)

	var answer struct {
		Driver           string `survey:"driver"`
		NamingConvention string `survey:"namingConvention"`
		Tag              string `survey:"tag"`
		Strict           bool   `survey:"strict,omitempty"`
	}

	if fi, _ := os.Stat(fileDest); fi != nil {
		log.Println(`configuration already exists`)
		return nil
	}

	if err := survey.Ask(questions, &answer); err != nil {
		return noInterruptError(err)
	}

	cfg := config.DefaultConfig()
	switch answer.Driver {
	case string(config.MySQL):
		cfg.Driver = config.MySQL
	case string(config.Postgres):
		cfg.Driver = config.Postgres
	case string(config.Sqlite):
		cfg.Driver = config.Sqlite
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

	w := bytes.NewBufferString("")
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer enc.Close()
	if err := enc.Encode(cfg); err != nil {
		return err
	}

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

	log.Println(`Creating ` + filename)
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
