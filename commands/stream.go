package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// StreamCommand is the interface for the stream command, which has subcommands for starting, stopping and backfilling streams.
type StreamCommand interface {
	internal.GenericCommand
}

// NewStreamCommand creates a new instance of the StreamCommand, which includes the start, stop and backfill subcommands.
func NewStreamCommand(start StreamStart, stop StreamStop, backfill StreamBackfill) StreamCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "stream",
		Short: "Allows users to start, stop or backfill a stream",
	}
	cmd.AddCommand(start.GetCommand())
	cmd.AddCommand(stop.GetCommand())
	cmd.AddCommand(backfill.GetCommand())

	return internal.NewGenericCommand(&cmd)
}
