package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// DowntimeStopCommand is a command to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
type DowntimeStopCommand interface {
	internal.GenericCommand
}

// NewDowntimeStopCommand creates a new instance of the DowntimeStopCommand, which allows users to stop downtime for a stream or a list of streams.
func NewDowntimeStopCommand(ds interfaces.DowntimeService) DowntimeStopCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "stop <key>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume",
		RunE:  ds.StopDowntime,
	}
	return internal.NewGenericCommand(&cmd)
}
