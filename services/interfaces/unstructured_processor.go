package interfaces

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// UnstructuredProcessor defines the interface for processing unstructured Kubernetes resources in the commands
// that executes logic on a resource list.
type UnstructuredProcessor interface {
	// Process takes a context and a namespaced name, retrieves the corresponding unstructured resource, and processes
	// it according to the command's logic. It returns the processed unstructured resource or an error if processing fails.
	Process(ctx context.Context, def types.NamespacedName) (*unstructured.Unstructured, error)
}
