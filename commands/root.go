package commands

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/internal"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// RootCommand is the main command required by Cobra to run the plugin. It is responsible for adding all subcommands and flags.
type RootCommand interface {
	internal.GenericCommand
}

// NewRootCommand creates a new RootCommand with the provided StreamCommand and DowntimeCommand as subcommands. It also adds the necessary flags for Kubernetes configuration.
func NewRootCommand(configFlags *genericclioptions.ConfigFlags, streamCommand StreamCommand, downtimeCommand DowntimeCommand) RootCommand { // coverage-ignore (trivial)
	rootCommand := &cobra.Command{
		Use: "kubectl-arcane",
	}
	rootCommand.AddCommand(streamCommand.GetCommand())
	rootCommand.AddCommand(downtimeCommand.GetCommand())

	configFlags.AddFlags(rootCommand.PersistentFlags())

	return internal.NewGenericCommand(rootCommand)
}
