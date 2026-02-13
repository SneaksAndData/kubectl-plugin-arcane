package services

import (
	"context"
	"fmt"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fieldManager = "kubectl-arcane"

// Ensure stream implements interfaces.StreamService
var _ interfaces.StreamService = (*stream)(nil)

// stream is a streamService that provides stream operations.
type stream struct {
	clientSet          *versioned.Clientset
	unstructuredClient client.Client
}

// NewStreamService creates a new instance of the stream, which provides stream operations.
func NewStreamService(clientSet *versioned.Clientset, unstructuredClient client.Client) interfaces.StreamService {
	return &stream{
		clientSet:          clientSet,
		unstructuredClient: unstructuredClient,
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
	return s.modifyStreamDefinition(ctx,
		parameters.Namespace,
		parameters.StreamClass,
		parameters.StreamId,
		streamapis.Running,
		func(def streamapis.Definition) error {
			return def.SetSuspended(false)
		},
		func(definition streamapis.Definition) bool {
			return definition.Suspended()
		},
	)
}

// Stop is a method that allows users to stop a stream, use the <key> parameter to identify the stream to stop
func (s *stream) Stop(ctx context.Context, parameters *models.StopParameters) error {
	return s.modifyStreamDefinition(ctx,
		parameters.Namespace,
		parameters.StreamClass,
		parameters.StreamId,
		streamapis.Suspended,
		func(def streamapis.Definition) error {
			return def.SetSuspended(true)
		},
		func(definition streamapis.Definition) bool {
			return !definition.Suspended()
		},
	)
}

func (s *stream) modifyStreamDefinition(ctx context.Context,
	namespace string,
	streamClass string,
	streamId string,
	expectedPhase streamapis.Phase,
	modifier func(streamapis.Definition) error,
	needModify func(streamapis.Definition) bool) error {

	sc, err := s.clientSet.StreamingV1().StreamClasses("").Get(ctx, streamClass, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching stream class: %w", err)
	}

	namespacedName := types.NamespacedName{Namespace: namespace, Name: streamId}
	streamDefinition, err := streamapis.GetStreamForClass(ctx, s.unstructuredClient, sc, namespacedName)
	if err != nil {
		return fmt.Errorf("error fetching stream definition: %w", err)
	}

	if !needModify(streamDefinition) {
		return errors.NewStatusNoOpError(expectedPhase, namespacedName)
	}

	err = modifier(streamDefinition)
	//err = streamDefinition.SetSuspended(false)
	if err != nil {
		return fmt.Errorf("error modifiing stream definition: %w", err)
	}

	err = s.unstructuredClient.Update(ctx, streamDefinition.ToUnstructured())
	if err != nil {
		return fmt.Errorf("error updating stream definition: %w", err)
	}

	return nil
}
