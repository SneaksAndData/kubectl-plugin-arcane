package interfaces

import (
	"context"

	"k8s.io/cli-runtime/pkg/printers"
)

// ExecutionQueue defines an interface for processing a queue of Kubernetes resources, applying a specified processor to
// each item.
type ExecutionQueue interface {

	// ProcessQueue processes items from the queue using the provided UnstructuredProcessor, printing results with the given ResourcePrinter,
	ProcessQueue(ctx context.Context, process UnstructuredProcessor, printer printers.ResourcePrinter, queuePublisher QueuePublisher) error
}
