package services

import (
	"sync"

	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ interfaces.ClientProvider = (*clientProvider)(nil)

// clientProvider is a struct that implements the ClientProvider interface.
// It lazily initializes both a typed clientset for the Arcane Operator's custom resources and ensures the client
// initialization order and caching of the clients.
type clientProvider struct {
	ConfigFlags *genericclioptions.ConfigFlags

	clientSetOnce sync.Once
	clientSet     *versioned.Clientset
	clientSetErr  error

	unstructuredOnce   sync.Once
	unstructuredClient client.Client
	unstructuredErr    error
}

func NewClientProvider(configFlags *genericclioptions.ConfigFlags) interfaces.ClientProvider {
	return &clientProvider{
		ConfigFlags: configFlags,
	}
}

func (cp *clientProvider) ProvideClientSet() (*versioned.Clientset, error) {
	cp.clientSetOnce.Do(func() {
		config, err := cp.ConfigFlags.ToRESTConfig()
		if err != nil {
			cp.clientSetErr = err
			return
		}
		cfg, err := versioned.NewForConfig(config)
		if err != nil {
			cp.clientSetErr = err
			return
		}
		cp.clientSet = cfg
	})
	return cp.clientSet, cp.clientSetErr
}

func (cp *clientProvider) ProvideUnstructuredClient() (client.Client, error) {
	cp.unstructuredOnce.Do(func() {
		config, err := cp.ConfigFlags.ToRESTConfig()
		if err != nil {
			cp.unstructuredErr = err
			return
		}
		c, err := client.New(config, client.Options{})
		if err != nil {
			cp.unstructuredErr = err
			return
		}
		cp.unstructuredClient = c
	})
	return cp.unstructuredClient, cp.unstructuredErr
}
