// Package providers package that contains function that provides necessary 3rd party dependencies for DI
package providers

import (
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

// ProvideConfigFlags creates a new instance of ConfigFlags and adds the necessary flags to the root command for Kubernetes configuration.
func ProvideConfigFlags(rootCommand commands.RootCommand) *genericclioptions.ConfigFlags { // coverage-ignore (trivial)
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(rootCommand.GetCommand().Flags())
	return configFlags
}

// ProvideRestConfig creates a new REST configuration from the provided ConfigFlags. This configuration is used to create Kubernetes clients for interacting with the cluster.
func ProvideRestConfig(configFlags *genericclioptions.ConfigFlags) (*rest.Config, error) { // coverage-ignore (trivial)
	restConfig, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return restConfig, nil
}

// ProvideClientSet creates a new instance of the versioned clientset using the provided REST configuration.
// This clientset is used to interact with the custom resources defined by the arcane-operator.
func ProvideClientSet(restConfig *rest.Config) (*versioned.Clientset, error) { // coverage-ignore (trivial)
	clientSet, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
