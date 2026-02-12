package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// DowntimeDeclareCommand is a command to declare downtime for a stream or a list of streams
type DowntimeDeclareCommand interface {
	internal.GenericCommand
}

// NewDowntimeDeclareCommand creates a new instance of the DowntimeDeclareCommand, which allows users to temporarily stop a stream or a list of streams.
func NewDowntimeDeclareCommand(ds interfaces.DowntimeService) DowntimeDeclareCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "declare <stream-class> <mask> <key>",
		Args:  cobra.ExactArgs(3),
		Short: "Begin downtime for a stream or a list of streams, use the <key> parameter to resume the stream(s) later",
		RunE:  ds.DeclareDowntime,
	}
	return internal.NewGenericCommand(&cmd)
}
