package services

import (
	"context"
	"os"
	"sync"

	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/util/workqueue"
)

type executionQueue struct {
	clientProvider cmdinterfaces.ClientProvider
}

func NewExecutionQueue(provider cmdinterfaces.ClientProvider) interfaces.ExecutionQueue {
	return &executionQueue{
		clientProvider: provider,
	}
}

func (s *executionQueue) ProcessQueue(ctx context.Context, process interfaces.UnstructuredProcessor, printer printers.ResourcePrinter, queuePublisher interfaces.QueuePublisher) error {
	rateLimiter := workqueue.DefaultTypedControllerRateLimiter[interfaces.QueueItem]()
	queue := workqueue.NewTypedRateLimitingQueue[interfaces.QueueItem](rateLimiter)
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

func (s *executionQueue) processObjects(ctx context.Context, queue interfaces.Queue, process interfaces.UnstructuredProcessor, printer printers.ResourcePrinter) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, shutdown := queue.Get()
			if shutdown {
				return
			}

			updated, hasUpdated, err := process.Process(ctx, item.Definition.NamespacedName(), item.Class)
			if err != nil {
				logging.LogError(item.Definition.ToUnstructured(), "modifying object, will retry later", err)
				queue.AddRateLimited(item)
				continue
			}

			if !hasUpdated { // coverage-ignore
				// If the processor indicates that there's no update needed
				queue.Forget(item)
				queue.Done(item)
				continue
			}

			unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
			if err != nil {
				logging.LogError(item.Definition.ToUnstructured(), "in constructing kubernetes client, will not retry", err)
				// If we can't get a client, there's no point in retrying, so we forget the item and move on.
				queue.Forget(item)
				queue.Done(item)
				continue
			}

			err = unstructuredClient.Update(ctx, updated)
			if err != nil {
				logging.LogError(item.Definition.ToUnstructured(), "updating client, will retry later", err)
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
