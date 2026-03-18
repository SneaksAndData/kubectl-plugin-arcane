package services

import (
	"context"

	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*DowntimeSummarizationProcessor)(nil)

type DowntimeSummarizationProcessor struct {
	reader      interfaces.UnstructuredReader
	streamClass string
	Summaries   map[string]int
}

func (s DowntimeSummarizationProcessor) Process(ctx context.Context, def types.NamespacedName) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, s.streamClass, def)
	if err != nil { // coverage-ignore
		return nil, false, err
	}

	labels := stream.GetLabels()

	if labels == nil {
		logging.LogError(stream, "has no labels, skipping", err)
		return nil, false, nil // Skip items that have no labels
	}

	s.Summaries[labels["arcane.sneaksanddata.com/downtime"]]++

	// We return nil here because we don't want to modify the original object, we just want to update our summaries
	return nil, false, nil
}
