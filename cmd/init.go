package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			questions := []*survey.Question{
				{
					Name: "naming",
					Prompt: &survey.Select{
						Message: "What is your naming convention:",
						Options: []string{},
						// Default: string(snakeCase),
					},
				},
			}

			answers := struct {
				NamingConvention string `survey:"naming"`
			}{}

			if err := survey.Ask(questions, &answers); err != nil {
				return err
			}

			log.Println(answers)

			return nil
		},
	}
}
