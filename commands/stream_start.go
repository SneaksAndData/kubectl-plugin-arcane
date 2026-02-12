package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
)

// StreamStart is a command that runs a stream start operation.
type StreamStart interface {
	internal.GenericCommand
}

// NewStreamStart creates a new instance of the StreamStart command, which runs a stream start operation.
func NewStreamStart(streamService interfaces.StreamService) StreamStart { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "stream <stream-class> <stream-id>",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream command",
		RunE: func(cmd *cobra.Command, args []string) error {
			startParameters, err := models.NewStartParameters(cmd, args)
			if err != nil {
				return err
			}
			return streamService.Start(cmd.Context(), startParameters)
		},
	}
	return internal.NewGenericCommand(&cmd)
}
