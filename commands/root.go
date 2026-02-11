package commands

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// RootCommand is the main command required by Cobra to run the plugin. It is responsible for adding all subcommands and flags.
type RootCommand interface {
	Execute() error
}

// NewRootCommand creates a new RootCommand with the provided StreamCommand and DowntimeCommand as subcommands. It also adds the necessary flags for Kubernetes configuration.
func NewRootCommand(streamCommand StreamCommand, downtimeCommand DowntimeCommand) RootCommand {
	rootCommand := &cobra.Command{
		Use: "kubectl-arcane",
	}
	rootCommand.AddCommand(streamCommand.GetCommand())
	rootCommand.AddCommand(downtimeCommand.GetCommand())

	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(rootCommand.Flags())
	return rootCommand
}
