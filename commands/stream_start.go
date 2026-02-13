package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StreamStart is a command that runs a stream start operation.
type StreamStart interface {
	internal.GenericCommand
}

// NewStreamStart creates a new instance of the StreamStart command, which runs a stream start operation.
func NewStreamStart(streamService interfaces.StreamService, configFlags *genericclioptions.ConfigFlags) StreamStart { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "start <stream-class> <stream-id>",
		Args:  cobra.ExactArgs(2),
		Short: "Start a stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			startParameters, err := models.NewStartParameters(cmd, args, configFlags)
			if err != nil {
				return err
			}
			return streamService.Start(cmd.Context(), startParameters)
		},
	}
	return internal.NewGenericCommand(&cmd)
}
