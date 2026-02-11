package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// 1. Initialize standard flags (kubeconfig, namespace, etc.)
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use: "kubectl-myplugin",
		Run: func(cmd *cobra.Command, args []string) {
			// 2. Generate REST config from flags
			config, _ := configFlags.ToRESTConfig()

			// 3. Create a Kubernetes client
			clientset, _ := kubernetes.NewForConfig(config)
			if clientset != nil {
				fmt.Println("hello")
			}

			// Your logic here (e.g., list pods)
		},
	}

	// 4. Register the flags with your command
	configFlags.AddFlags(cmd.Flags())
	cobra.CheckErr(cmd.Execute())
}

func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(commands.CobraCommand)),
		fx.ResultTags(`group:"commands"`),
	)
}
