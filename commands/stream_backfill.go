package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// StreamBackfill is a command that runs a stream backfill operation.
type StreamBackfill interface {
	internal.GenericCommand
}

// NewStreamBackfill creates a new instance of the StreamBackfill command, which runs a stream backfill operation.
func NewStreamBackfill(operator interfaces.StreamService) StreamBackfill {
	cmd := cobra.Command{
		Use:   "backfill <stream-class> <stream-id> [--wait]",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream command",
		RunE:  operator.Backfill,
	}
	return internal.NewGenericCommand(&cmd)
}
