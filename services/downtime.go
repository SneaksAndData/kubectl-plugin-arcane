package services

import (
	"context"
	"os"
	"sync"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/filter"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/util/workqueue"
)

// Ensure downtime implements cmdinterfaces.DowntimeService
var _ cmdinterfaces.DowntimeService = (*downtime)(nil)

// Queue is a typed rate-limiting work queue for unstructured objects.
type Queue = workqueue.TypedRateLimitingInterface[streamapis.Definition]

// UnstructuredObjectFilter is a function that filters unstructured objects based on custom criteria.
type UnstructuredObjectFilter func(item streamapis.Definition) bool

// downtime is a service that provides downtime operations.
type downtime struct {
	clientProvider cmdinterfaces.ClientProvider
	factory        *DowntimeProcessorFactory
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService(clientProvider cmdinterfaces.ClientProvider, factory *DowntimeProcessorFactory) cmdinterfaces.DowntimeService {
	return &downtime{
		clientProvider: clientProvider,
		factory:        factory,
	}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(ctx context.Context, parameters *models.DowntimeDeclareParameters) error {
	f := filter.NewUnsuspendedByNamePrefix(parameters.Prefix)
	publisher := NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, parameters.Namespace, f)
	return s.runWithQueue(
		ctx,
		s.factory.DowntimeDeclareProcessor(parameters),
		Printer("suspended"),
		publisher,
	)
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(ctx context.Context, parameters *models.DowntimeStopParameters) error {
	f := filter.NewByDowntimeKey(parameters.DowntimeKey)
	publisher := NewStreamClassMembersPublisher(s.clientProvider, parameters.StreamClass, "", f)
	return s.runWithQueue(ctx,
		s.factory.DowntimeStopProcessor(parameters),
		Printer("started"),
		publisher,
	)
}

func (s *downtime) runWithQueue(ctx context.Context, process interfaces.UnstructuredProcessor, printer printers.ResourcePrinter, queuePublisher interfaces.QueuePublisher) error {
	rateLimiter := workqueue.DefaultTypedControllerRateLimiter[streamapis.Definition]()
	queue := workqueue.NewTypedRateLimitingQueue[streamapis.Definition](rateLimiter)
	defer queue.ShutDown()
	var wg sync.WaitGroup

	wg.Go(func() {
		s.processObjects(ctx, queue, process, printer)
	})

	err := queuePublisher.PublishStreamDefinitions(ctx, queue)
	if err != nil { // coverage-ignore
		return err
	}

	queue.ShutDownWithDrain()
	wg.Wait()
	return nil
}

func (s *downtime) processObjects(ctx context.Context, queue Queue, process interfaces.UnstructuredProcessor, printer printers.ResourcePrinter) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, shutdown := queue.Get()
			if shutdown {
				return
			}

			updated, err := process.Process(ctx, item.NamespacedName())
			if err != nil {
				logError(item.ToUnstructured(), "modifying object, will retry later", err)
				queue.AddRateLimited(item)
				continue
			}

			unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
			if err != nil {
				logError(item.ToUnstructured(), "in constructing kubernetes client, will not retry", err)
				// If we can't get a client, there's no point in retrying, so we forget the item and move on.
				queue.Forget(item)
				queue.Done(item)
				continue
			}

			err = unstructuredClient.Update(ctx, updated)
			if err != nil {
				logError(item.ToUnstructured(), "updating client, will retry later", err)
				queue.AddRateLimited(item)
				continue
			}

			queue.Forget(item)
			queue.Done(item)
			err = printer.PrintObj(updated, os.Stdout)
			if err != nil {
				// If we can't print, we still consider the item processed successfully, so we forget it and move on.
				continue
			}
		}
	}
}
