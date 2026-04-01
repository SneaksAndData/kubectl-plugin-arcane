package services

import (
	"context"
	"time"

	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*downtimeDeclareProcessor)(nil)

type downtimeDeclareProcessor struct {
	key    string
	reader interfaces.UnstructuredReader
}

func (s *downtimeDeclareProcessor) Process(ctx context.Context, def types.NamespacedName, class *v1.StreamClass) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, class, def)
	if err != nil {
		return nil, false, err
	}

	labels := stream.GetLabels()

	if existingKey, exists := labels[interfaces.DowntimeLabelKey]; exists && existingKey != s.key {
		logging.LogError(stream, "already has a different downtime key", err)
		return nil, false, nil // Skip items that already have a different downtime key
	}

	if labels == nil {
		labels = make(map[string]string)
	}

	labels[interfaces.DowntimeLabelKey] = s.key
	stream.SetLabels(labels)

	annotations := stream.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[interfaces.DowntimeBeginAnnotationKey] = time.Now().UTC().Format(time.RFC3339)
	stream.SetAnnotations(annotations)

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
