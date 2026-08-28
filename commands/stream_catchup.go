package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StreamCatchup is a command that runs a stream catchup operation.
type StreamCatchup interface {
	internal.GenericCommand
}

// NewStreamCatchup creates a new instance of the StreamCatchup command, which runs a stream catchup operation.
func NewStreamCatchup(catchupService interfaces.BackfillService, configFlags *genericclioptions.ConfigFlags) StreamCatchup { // coverage-ignore
	var overrides []string
	overrides = append(overrides, ".spec.backfillBehavior=Merge")
	cmd := cobra.Command{
		Use:   "catchup <stream-class> <stream-id> [--wait]",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream in the backfill mode with the backfill behavior set to Merge",
		RunE: func(cmd *cobra.Command, args []string) error {
			parameters, err := models.NewBackfillParameters(cmd, args, configFlags, &overrides)
			if err != nil {
				return err
			}
			return catchupService.Backfill(cmd.Context(), parameters)
		},
	}
	command := catchupCommand{
		Command:   &cmd,
		overrides: &overrides,
	}
	cmd.Flags().Bool("wait", false, "Wait for catchup command to complete")
	cmd.Flags().StringArrayVarP(command.overrides, "override", "o", []string{}, "Override additional spec values (format: key=value)")
	return &command
}

type catchupCommand struct {
	*cobra.Command
	overrides *[]string
}

func (b *catchupCommand) GetCommand() *cobra.Command { // coverage-ignore (trivial)
	return b.Command
}
