package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "my-app",
		Short: "A brief description",
		Run: func(cmd *cobra.Command, args []string) {
			// Code that runs when no subcommand is provided
		},
	}
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(rootCmd.PersistentFlags())

	fx.Supply(configFlags)
}
