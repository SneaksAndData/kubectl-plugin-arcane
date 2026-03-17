package interfaces

import (
	"context"

	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"k8s.io/client-go/util/workqueue"
)

// Queue is a typed rate-limiting work queue for unstructured objects.
type Queue = workqueue.TypedRateLimitingInterface[stream.Definition]

type QueuePublisher interface {
	// PublishStreamDefinitions retrieves a list of objects based on the provided parameters.
	PublishStreamDefinitions(ctx context.Context, target Queue) error
}
