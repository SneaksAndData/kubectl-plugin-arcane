package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// StreamStart is a command that runs a stream start operation.
type StreamStart interface {
	internal.GenericCommand
}

// NewStreamStart creates a new instance of the StreamStart command, which runs a stream start operation.
func NewStreamStart(operator interfaces.StreamService) StreamStart {
	cmd := cobra.Command{
		Use:   "stream <stream-class> <stream-id>",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream command",
		RunE:  operator.Start,
	}
	return internal.NewGenericCommand(&cmd)
}
