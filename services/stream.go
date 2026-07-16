package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SneaksAndData/arcane-operator/services/controllers/contracts"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/errors"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/printers"
)

const fieldManager = "kubectl-arcane"

// Ensure stream implements interfaces.StreamService
var _ interfaces.StreamService = (*stream)(nil)

// stream is a service that provides stream operations.
type stream struct {
	clientProvider interfaces.ClientProvider
}

// NewStreamService creates a new instance of the stream, which provides stream operations.
func NewStreamService(clientProvider interfaces.ClientProvider) interfaces.StreamService {
	return &stream{
		clientProvider: clientProvider,
	}
}

// Start is a method that allows users to start a stream, use the <key> parameter to identify the stream to start
func (s *stream) Start(ctx context.Context, parameters *models.StartParameters) error {
	printer := logging.Printer("started")
	return wait.PollUntilContextCancel(ctx, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		err = s.modifyStreamDefinition(ctx,
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
			printer,
		)
		if err == nil || apierrors.IsConflict(err) {
			return true, nil
		}
		return false, err
	})
}

// Stop is a method that allows users to stop a stream, use the <key> parameter to identify the stream to stop
func (s *stream) Stop(ctx context.Context, parameters *models.StopParameters) error {
	printer := logging.Printer("stopped")
	return wait.PollUntilContextCancel(ctx, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		err = s.modifyStreamDefinition(ctx,
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
			printer,
		)
		if err == nil || apierrors.IsConflict(err) {
			return true, nil
		}
		return false, err
	})
}

func (s *stream) modifyStreamDefinition(ctx context.Context,
	namespace string,
	streamClass string,
	streamId string,
	expectedPhase streamapis.Phase,
	modifier func(streamapis.Definition) error,
	needModify func(streamapis.Definition) bool, printer printers.ResourcePrinter) error {

	clientSet, err := s.clientProvider.ProvideClientSet()
	if err != nil {
		return fmt.Errorf("error providing client set: %w", err)
	}
	sc, err := clientSet.StreamingV1().StreamClasses("").Get(ctx, streamClass, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching stream class: %w", err)
	}

	namespacedName := types.NamespacedName{Namespace: namespace, Name: streamId}
	unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
	if err != nil {
		return fmt.Errorf("error providing unstructured client: %w", err)
	}

	streamDefinition, err := streamapis.GetStreamForClass(ctx, unstructuredClient, sc, namespacedName, contracts.FromUnstructured)
	if err != nil {
		return fmt.Errorf("error fetching stream definition: %w", err)
	}

	if !needModify(streamDefinition) {
		return errors.NewStatusNoOpError(expectedPhase, namespacedName)
	}

	err = modifier(streamDefinition)
	if err != nil {
		return fmt.Errorf("error modifiing stream definition: %w", err)
	}

	err = unstructuredClient.Update(ctx, streamDefinition.ToUnstructured())
	if err != nil {
		return fmt.Errorf("error updating stream definition: %w", err)
	}

	err = printer.PrintObj(streamDefinition.ToUnstructured(), os.Stdout)
	if err != nil {
		return fmt.Errorf("error printing stream definition: %w", err)
	}
	return nil
}
