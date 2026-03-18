package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
)

// DowntimeListCommand is a command to list active downtime keys in the cluster, optionally filtered by stream class
type DowntimeListCommand interface {
	internal.GenericCommand
}

// NewDowntimeListCommand creates a new instance of the DowntimeListCommand, which allows users to stop downtime for a stream or a list of streams.
func NewDowntimeListCommand(ds interfaces.DowntimeService) DowntimeListCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "list",
		Args:  cobra.NoArgs,
		Short: "List of active downtime keys in the cluster, optionally filtered by stream class",
		RunE: func(cmd *cobra.Command, args []string) error {
			parameters, err := models.NewDowntimeListParameters(cmd)
			if err != nil {
				return err
			}
			return ds.ListDowntimes(cmd.Context(), parameters)
		},
	}

	// add --stream-class flag so callers can filter by stream class without positional args
	cmd.Flags().String("stream-class", "", "Filter by stream class")

	return internal.NewGenericCommand(&cmd)
}
