package services

import (
	"context"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*downtimeDeclareProcessor)(nil)

type downtimeDeclareProcessor struct {
	key         string
	reader      interfaces.UnstructuredReader
	streamClass string
}

func (s *downtimeDeclareProcessor) Process(ctx context.Context, def types.NamespacedName) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, s.streamClass, def)
	if err != nil {
		return nil, false, err
	}

	labels := stream.GetLabels()

	if existingKey, exists := labels["arcane.sneaksanddata.com/downtime"]; exists && existingKey != s.key {
		logging.LogError(stream, "already has a different downtime key", err)
		return nil, false, nil // Skip items that already have a different downtime key
	}

	if labels == nil {
		labels = make(map[string]string)
	}

	labels["arcane.sneaksanddata.com/downtime"] = s.key
	stream.SetLabels(labels)

	definition, err := streamapis.FromUnstructured(stream)
	if err != nil {
		return nil, false, err
	}
	err = definition.SetSuspended(true)
	if err != nil {
		return nil, false, err
	}
	return definition.ToUnstructured(), true, nil
}
