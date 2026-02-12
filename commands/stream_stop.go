package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
)

// StreamStop is a command that runs a stream stop operation.
type StreamStop interface {
	internal.GenericCommand
}

// NewStreamStop creates a new instance of the StreamStop command, which runs a stream stop operation.
func NewStreamStop(streamService interfaces.StreamService) StreamStop { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "stop <stream-class> <stream-id>",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream command",
		RunE: func(cmd *cobra.Command, args []string) error {
			startParameters, err := models.NewStopParameters(cmd, args)
			if err != nil {
				return err
			}
			return streamService.Stop(cmd.Context(), startParameters)
		},
	}
	return internal.NewGenericCommand(&cmd)
}
