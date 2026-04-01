package commands

import (
	"fmt"

	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
)

// VersionCommand is a command that prints the plugin version.
type VersionCommand interface {
	internal.GenericCommand
}

// NewVersionCommand creates a version command that prints the injected build version.
func NewVersionCommand(version string, buildNumber string) VersionCommand { // coverage-ignore (trivial)
	cmd := cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Print plugin version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println(fmt.Sprintf("kuberctl plugin arcane. Version: %s, Build Number: %s", version, buildNumber))
			return nil
		},
	}

	return internal.NewGenericCommand(&cmd)
}
