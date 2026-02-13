package services

import (
	"context"
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"sync"
)

// Ensure downtime implements interfaces.DowntimeService
var _ interfaces.DowntimeService = (*downtime)(nil)

// Queue is a typed rate-limiting work queue for unstructured objects.
type Queue = workqueue.TypedRateLimitingInterface[*unstructured.Unstructured]

// UnstructuredProcessor is a function that processes an unstructured object and returns an updated unstructured object or an error.
type UnstructuredProcessor func(item *unstructured.Unstructured) (*unstructured.Unstructured, error)

// downtime is a service that provides downtime operations.
type downtime struct {
	streamService      interfaces.StreamService
	clientSet          *versioned.Clientset
	unstructuredClient client.Client
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService(clientSet *versioned.Clientset, unstructuredClient client.Client) interfaces.DowntimeService {
	return &downtime{
		clientSet:          clientSet,
		unstructuredClient: unstructuredClient,
	}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(ctx context.Context, parameters *models.DowntimeDeclareParameters) error {
	return s.runWithQueue(ctx, parameters.StreamClass, parameters.Namespace, parameters.Prefix, setDowntimeForStream(parameters.DowntimeKey))
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(ctx context.Context, parameters *models.DowntimeStopParameters) error {
	return s.runWithQueue(ctx, parameters.StreamClass, parameters.Namespace, parameters.Prefix, unsetDowntimeForStream(parameters.DowntimeKey))
}

func setDowntimeForStream(key string) UnstructuredProcessor {
	return func(item *unstructured.Unstructured) (*unstructured.Unstructured, error) {
		labels := item.GetLabels()

		if labels == nil {
			labels = make(map[string]string)
		}

		// Skip if already has a downtime key that's different
		if existingKey, exists := labels["arcane.sneaksanddata.com/downtime"]; exists && existingKey != key {
			return nil, nil
		}

		labels["arcane.sneaksanddata.com/downtime"] = key
		item.SetLabels(labels)

		definition, err := streamapis.FromUnstructured(item)
		if err != nil {
			return nil, err
		}
		err = definition.SetSuspended(true)
		if err != nil {
			return nil, err
		}

		return definition.ToUnstructured(), nil
	}
}

func unsetDowntimeForStream(key string) UnstructuredProcessor {
	return func(item *unstructured.Unstructured) (*unstructured.Unstructured, error) {
		labels := item.GetLabels()

		if labels["arcane.sneaksanddata.com/downtime"] != key {
			return nil, nil // Skip items that don't match the downtime key
		}

		delete(labels, "arcane.sneaksanddata.com/downtime")
		item.SetLabels(labels)

		definition, err := streamapis.FromUnstructured(item)
		if err != nil {
			return nil, err
		}

		err = definition.SetSuspended(false)
		if err != nil {
			return nil, err
		}

		return definition.ToUnstructured(), nil
	}
}

func (s *downtime) runWithQueue(ctx context.Context, streamClass string, namespace string, prefix string, process UnstructuredProcessor) error {
	rateLimiter := workqueue.DefaultTypedControllerRateLimiter[*unstructured.Unstructured]()
	queue := workqueue.NewTypedRateLimitingQueue[*unstructured.Unstructured](rateLimiter)
	defer queue.ShutDown()
	var wg sync.WaitGroup

	wg.Go(func() {
		s.processObjects(ctx, queue, process)
	})

	err := s.getObjectsList(ctx, streamClass, namespace, prefix, queue)
	if err != nil {
		return err
	}

	queue.ShutDownWithDrain()
	wg.Wait()
	return nil
}

func (s *downtime) getObjectsList(ctx context.Context, streamClass string, namespace string, prefix string, queue Queue) error {
	sc, err := s.clientSet.StreamingV1().StreamClasses("").Get(ctx, streamClass, metav1.GetOptions{})
	if err != nil {
		return err
	}

	gvk := sc.TargetResourceGvk()

	streamList := &unstructured.UnstructuredList{}
	streamList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind + "List",
	})

	err = s.unstructuredClient.List(ctx, streamList, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	for _, item := range streamList.Items {
		if !strings.HasPrefix(item.GetName(), prefix) {
			continue
		}
		queue.Add(&item)
	}

	return nil
}

func (s *downtime) processObjects(ctx context.Context, queue Queue, process UnstructuredProcessor) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, shutdown := queue.Get()
			if shutdown {
				return
			}

			itemCopy := item.DeepCopy()
			updated, err := process(itemCopy)
			if err != nil {
				queue.AddRateLimited(item)
				continue
			}

			if updated == nil {
				queue.Forget(item)
				queue.Done(item)
				continue
			}

			err = s.unstructuredClient.Update(ctx, updated)
			if err != nil {
				queue.AddRateLimited(item)
				continue
			}

			queue.Forget(item)
			queue.Done(item)
		}
	}

}
