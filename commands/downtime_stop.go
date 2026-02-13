package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			parameters, err := models.NewDowntimeStopParameters(cmd, args)
			if err != nil {
				return err
			}
			return ds.StopDowntime(cmd.Context(), parameters)
		},
	}
	return internal.NewGenericCommand(&cmd)
}
