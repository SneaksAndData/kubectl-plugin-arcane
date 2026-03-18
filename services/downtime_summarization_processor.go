package services

import (
	"context"
	"fmt"

	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*DowntimeSummarizationProcessor)(nil)

type DowntimeSummarizationProcessor struct {
	reader  interfaces.UnstructuredReader
	Summary map[string][]string
}

func NewDowntimeSummarizationProcessor(reader interfaces.UnstructuredReader) *DowntimeSummarizationProcessor {
	return &DowntimeSummarizationProcessor{
		reader:  reader,
		Summary: make(map[string][]string),
	}
}

func (s DowntimeSummarizationProcessor) Process(ctx context.Context, def types.NamespacedName, class *v1.StreamClass) (*unstructured.Unstructured, bool, error) {
	stream, err := s.reader.Read(ctx, class, def)
	if err != nil { // coverage-ignore
		return nil, false, err
	}

	labels := stream.GetLabels()

	if labels == nil { // coverage-ignore
		return nil, false, nil
	}

	label := labels[interfaces.DowntimeAnnotationKey]
	s.Summary[label] = append(s.Summary[label], fmt.Sprintf("%s/%s", stream.GetNamespace(), stream.GetName()))

	// We return nil here because we don't want to modify the original object, we just want to update our summaries
	return nil, false, nil
}
