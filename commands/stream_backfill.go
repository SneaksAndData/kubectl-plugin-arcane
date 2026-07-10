package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StreamBackfill is a command that runs a stream backfill operation.
type StreamBackfill interface {
	internal.GenericCommand
	GetOverrides() []string
}

// NewStreamBackfill creates a new instance of the StreamBackfill command, which runs a stream backfill operation.
func NewStreamBackfill(streamService interfaces.StreamService, configFlags *genericclioptions.ConfigFlags) StreamBackfill { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "backfill <stream-class> <stream-id> [--wait]",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream in backfill mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			parameters, err := models.NewBackfillParameters(cmd, args, configFlags)
			if err != nil {
				return err
			}
			return streamService.Backfill(cmd.Context(), parameters)
		},
	}
	command := backfillCommand{
		Command:   &cmd,
		overrides: []string{},
	}
	cmd.Flags().Bool("wait", false, "Wait for backfill command to complete")
	cmd.Flags().StringArrayVarP(&command.overrides, "override", "o", []string{}, "Override spec values (format: key=value)")
	return &command
}

type backfillCommand struct {
	*cobra.Command
	overrides []string
}

func (b *backfillCommand) GetCommand() *cobra.Command {
	return b.Command
}

func (b *backfillCommand) GetOverrides() []string {
	return b.overrides
}
