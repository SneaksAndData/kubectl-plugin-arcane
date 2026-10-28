package commands

import (
	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StreamMigrate is a command that runs a stream migration operation.
type StreamMigrate interface {
	internal.GenericCommand
}

// NewStreamMigrate creates a new instance of the StreamMigrate command, which runs a stream migration operation.
// This command does not support overrides.
func NewStreamMigrate(backfillService interfaces.BackfillService, configFlags *genericclioptions.ConfigFlags) StreamMigrate { // coverage-ignore
	cmd := cobra.Command{
		Use:   "migrate <stream-class> <stream-id> [--wait]",
		Args:  cobra.ExactArgs(2),
		Short: "Run a stream in migration mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			parameters, err := models.NewBackfillParameters(cmd, args, configFlags, &[]string{}, v1.BackfillTypeMigration)
			if err != nil {
				return err
			}
			return backfillService.Backfill(cmd.Context(), parameters)
		},
	}
	command := backfillCommand{
		Command:   &cmd,
		overrides: &[]string{},
	}
	cmd.Flags().Bool("wait", false, "Wait for migration command to complete")
	return &command
}
