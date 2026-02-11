package providers

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

// ProvideConfigFlags creates a new instance of ConfigFlags and adds the necessary flags to the root command for Kubernetes configuration.
func ProvideConfigFlags(rootCommand commands.RootCommand) *genericclioptions.ConfigFlags {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(rootCommand.GetCommand().Flags())
	return configFlags
}

// ProvideRestConfig creates a new REST configuration from the provided ConfigFlags. This configuration is used to create Kubernetes clients for interacting with the cluster.
func ProvideRestConfig(configFlags *genericclioptions.ConfigFlags) (*rest.Config, error) {
	restConfig, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return restConfig, nil
}
