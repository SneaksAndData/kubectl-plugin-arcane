package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// StreamStop is a command that runs a stream stop operation.
type StreamStop interface {
	internal.GenericCommand
}

// NewStreamStop creates a new instance of the StreamStop command, which runs a stream stop operation.
func NewStreamStop(operator interfaces.StreamService) StreamStop {
	cmd := cobra.Command{
		Use:   "stop <stream-class> <stream-id>",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream command",
		RunE:  operator.Stop,
	}
	return internal.NewGenericCommand(&cmd)
}
