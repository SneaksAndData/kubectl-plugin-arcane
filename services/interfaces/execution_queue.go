package interfaces

import (
	"context"

	"k8s.io/cli-runtime/pkg/printers"
)

type ExecutionQueue interface {
	ProcessQueue(ctx context.Context, process UnstructuredProcessor, printer printers.ResourcePrinter, queuePublisher QueuePublisher) error
}
