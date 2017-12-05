package hook

import (
	"fmt"

	"github.com/sensu/sensu-go/cli"
	"github.com/spf13/cobra"
)

// UpdateCommand adds command that allows user to create new hooks
func UpdateCommand(cli *cli.SensuCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update NAME",
		Short:        "update hooks",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print ot usage if we do not receive one argument
			if len(args) != 1 {
				cmd.Help()
				return nil
			}

			// Fetch hooks from API
			hookID := args[0]
			hook, err := cli.Client.FetchHook(hookID)
			if err != nil {
				return err
			}

			// Administer questionnaire
			opts := newHookOpts()
			opts.withHook(hook)
			if err := opts.administerQuestionnaire(true); err != nil {
				return err
			}

			// Apply given arguments to hook
			opts.Copy(hook)

			if err := hook.Validate(); err != nil {
				return err
			}

			//
			// TODO:
			//
			// Current validation is a bit too laissez faire. For usability we should
			// determine whether there are assets / handlers / mutators associated w/
			// the hook and warn the user if they do not exist yet.
			if err := cli.Client.UpdateHook(hook); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OK")
			return nil
		},
	}

	return cmd
}