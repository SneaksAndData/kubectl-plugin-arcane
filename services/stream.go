package services

import (
	"context"
	"fmt"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const fieldManager = "kubectl-arcane"

// Ensure stream implements interfaces.StreamService
var _ interfaces.StreamService = (*stream)(nil)

// stream is a service that provides stream operations.
type stream struct {
	clientSet *versioned.Clientset
}

// NewStreamService creates a new instance of the stream, which provides stream operations.
func NewStreamService(clientSet *versioned.Clientset) interfaces.StreamService {
	return &stream{
		clientSet: clientSet,
	}
}

// Backfill is a method that allows users to run a stream backfill operation
func (s *stream) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {
	bfr, err := s.clientSet.
		StreamingV1().
		BackfillRequests(parameters.Namespace).
		Create(ctx, parameters.ToBackfillRequest(), metav1.CreateOptions{
			DryRun:          parameters.DryRun,
			FieldManager:    fieldManager,
			FieldValidation: "Strict",
		})
	if err != nil {
		return fmt.Errorf("error creating backfill request: %w", err)
	}

	if !parameters.Wait {
		return nil
	}
	watch, err := s.clientSet.StreamingV1().BackfillRequests(parameters.Namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector:   "metadata.name=" + bfr.Name,
		Watch:           true,
		ResourceVersion: bfr.ResourceVersion,
	})
	if err != nil {
		return fmt.Errorf("error watching backfill request: %w", err)
	}
	defer watch.Stop()

	for {
		select {
		case event, ok := <-watch.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}
			bfr, ok := event.Object.(*v1.BackfillRequest)
			if !ok {
				return fmt.Errorf("unexpected object type: %T", event.Object)
			}

			if bfr.Spec.Completed {
				return nil
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Start is a method that allows users to start a stream, use the <key> parameter to identify the stream to start
func (s *stream) Start(ctx context.Context, parameters *models.StartParameters) error {
	panic("not implemented")
}

// Stop is a method that allows users to stop a stream, use the <key> parameter to identify the stream to stop
func (s *stream) Stop(ctx context.Context, parameters *models.StopParameters) error {
	panic("not implemented")
}
