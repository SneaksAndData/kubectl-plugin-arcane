package interfaces

import (
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClientProvider is an interface that provides methods to obtain Kubernetes clients, including both typed and unstructured clients.
// This interface is intended to lazily initialize clients when they are first needed.
// The methods of this client should never be called in class constructors.
type ClientProvider interface {
	// ProvideClientSet returns a typed clientset for the Arcane Operator's custom resources.
	ProvideClientSet() (*versioned.Clientset, error)

	// ProvideUnstructuredClient returns a controller-runtime client that can be used for unstructured operations.
	ProvideUnstructuredClient() (client.Client, error)
}
