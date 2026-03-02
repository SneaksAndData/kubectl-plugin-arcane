package interfaces

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// UnstructuredReader defines the interface for reading unstructured Kubernetes resources.
type UnstructuredReader interface {
	// Read retrieves an unstructured Kubernetes resource based on the provided namespaced name.
	Read(ctx context.Context, streamClass string, name types.NamespacedName) (*unstructured.Unstructured, error)
}
