package services

import (
	"context"

	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*downtimeStopProcessor)(nil)

type downtimeStopProcessor struct {
	key    string
	reader interfaces.UnstructuredReader
}

func (s downtimeStopProcessor) Process(ctx context.Context, def types.NamespacedName, class *v1.StreamClass) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, class, def)
	if err != nil { // coverage-ignore
		return nil, false, err
	}

	labels := stream.GetLabels()

	if labels[interfaces.DowntimeLabelKey] != s.key {
		logging.LogError(stream, "has a different downtime key, skipping", err)
		return nil, false, nil // Skip items that don't match the downtime key
	}

	delete(labels, interfaces.DowntimeLabelKey)
	delete(labels, interfaces.DowntimeBeginLabelKey)
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
