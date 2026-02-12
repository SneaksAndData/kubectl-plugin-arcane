package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// DowntimeCommand is the interface for the downtime command, which allows users to temporarily stop a stream or a list of streams.
type DowntimeCommand interface {
	internal.GenericCommand
}

// NewDowntimeCommand creates a new instance of the DowntimeCommand, which includes the declare and stop subcommands.
func NewDowntimeCommand(command DowntimeDeclareCommand, stopCommand DowntimeStopCommand) DowntimeCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "downtime",
		Short: "Temporarily stop or start a stream or a list of streams",
	}
	cmd.AddCommand(command.GetCommand())
	cmd.AddCommand(stopCommand.GetCommand())
	return internal.NewGenericCommand(&cmd)
}
