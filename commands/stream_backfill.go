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
	cmd.Flags().Bool("wait", false, "Wait for backfill command to complete")
	return internal.NewGenericCommand(&cmd)
}
