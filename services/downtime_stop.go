package services

import (
	"context"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*downtimeStopProcessor)(nil)

type downtimeStopProcessor struct {
	key         string
	reader      interfaces.UnstructuredReader
	streamClass string
}

func (s downtimeStopProcessor) Process(ctx context.Context, def types.NamespacedName) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, s.streamClass, def)
	if err != nil { // coverage-ignore
		return nil, false, err
	}

	labels := stream.GetLabels()

	if labels["arcane.sneaksanddata.com/downtime"] != s.key {
		logging.LogError(stream, "has a different downtime key, skipping", err)
		return nil, false, nil // Skip items that don't match the downtime key
	}

	delete(labels, "arcane.sneaksanddata.com/downtime")
	stream.SetLabels(labels)

	definition, err := streamapis.FromUnstructured(stream)
	if err != nil { // coverage-ignore
		return nil, false, err
	}
	err = definition.SetSuspended(false)
	if err != nil { // coverage-ignore
		return nil, false, err
	}
	return definition.ToUnstructured(), true, nil
}
