package services

import (
	"context"

	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredReader = (*unstructuredReader)(nil)

type unstructuredReader struct {
	clientProvider cmdinterfaces.ClientProvider
}

func NewUnstructuredReader(clientProvider cmdinterfaces.ClientProvider) interfaces.UnstructuredReader {
	return &unstructuredReader{
		clientProvider: clientProvider,
	}
}

func (s *unstructuredReader) Read(ctx context.Context, streamClass *v1.StreamClass, name types.NamespacedName) (*unstructured.Unstructured, error) {
	gvk := streamClass.TargetResourceGvk()
	stream := &unstructured.Unstructured{}
	stream.SetGroupVersionKind(gvk)
	stream.SetNamespace(name.Namespace)
	stream.SetName(name.Name)

	unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
	if err != nil {
		return nil, err
	}
	err = unstructuredClient.Get(ctx, name, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
