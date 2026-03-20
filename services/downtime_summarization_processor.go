package services

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.UnstructuredProcessor = (*DowntimeSummarizationProcessor)(nil)

type DowntimeSummarizationProcessor struct {
	reader    interfaces.UnstructuredReader
	Summary   map[string][]string
	Durations map[string]time.Time
}

func NewDowntimeSummarizationProcessor(reader interfaces.UnstructuredReader) *DowntimeSummarizationProcessor {
	return &DowntimeSummarizationProcessor{
		reader:    reader,
		Summary:   make(map[string][]string),
		Durations: make(map[string]time.Time),
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

	label := labels[interfaces.DowntimeLabelKey]

	var ms time.Time
	annotations := stream.GetAnnotations()
	if annotations != nil {
		startDate := annotations[interfaces.DowntimeBeginAnnotationKey]
		ms, err = time.ParseInLocation(time.RFC3339, startDate, time.UTC)
		if err != nil {
			logging.LogError(stream, "to parse downtime start date for stream, skipping", err)
			ms = time.Now().UTC()
		}
	} else {
		logging.LogError(stream, "to parse downtime start date for stream, skipping", err)
		ms = time.Now().UTC()
	}
	// Go 1.17+:
	streamId := fmt.Sprintf("%s/%s", stream.GetNamespace(), stream.GetName())
	s.Summary[label] = append(s.Summary[label], streamId)

	// We want to keep the earliest downtime start time for each key
	if prev, ok := s.Durations[label]; !ok || ms.Before(prev) {
		s.Durations[label] = ms
	}

	// We return nil here because we don't want to modify the original object, we just want to update our summaries
	return nil, false, nil
}
