package services

import (
	"context"
	"fmt"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/util/workqueue"
	"os"
	"strings"
	"sync"
)

// Ensure downtime implements interfaces.DowntimeService
var _ interfaces.DowntimeService = (*downtime)(nil)

// Queue is a typed rate-limiting work queue for unstructured objects.
type Queue = workqueue.TypedRateLimitingInterface[streamapis.Definition]

// UnstructuredProcessor is a function that processes an unstructured object and returns an updated unstructured object or an error.
type UnstructuredProcessor func(def streamapis.Definition) error

// UnstructuredObjectFilter is a function that filters unstructured objects based on custom criteria.
type UnstructuredObjectFilter func(item streamapis.Definition) bool

// downtime is a service that provides downtime operations.
type downtime struct {
	clientProvider interfaces.ClientProvider
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService(clientProvider interfaces.ClientProvider) interfaces.DowntimeService {
	return &downtime{
		clientProvider: clientProvider,
	}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(ctx context.Context, parameters *models.DowntimeDeclareParameters) error {
	return s.runWithQueue(
		ctx,
		parameters.StreamClass,
		filterByNamePrefix(parameters.Prefix),
		setDowntimeForStream(parameters.DowntimeKey),
		Printer("suspended"),
	)
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(ctx context.Context, parameters *models.DowntimeStopParameters) error {
	return s.runWithQueue(ctx,
		parameters.StreamClass,
		filterByDowntimeKey(parameters.DowntimeKey),
		unsetDowntimeForStream(parameters.DowntimeKey),
		Printer("started"),
	)
}

func setDowntimeForStream(key string) UnstructuredProcessor {
	return func(def streamapis.Definition) error {
		item := def.ToUnstructured()
		labels := item.GetLabels()

		if labels == nil {
			labels = make(map[string]string)
		}

		// Skip if already has a downtime key that's different
		if existingKey, exists := labels["arcane.sneaksanddata.com/downtime"]; exists && existingKey != key {
			return nil
		}

		labels["arcane.sneaksanddata.com/downtime"] = key
		item.SetLabels(labels)

		definition, err := streamapis.FromUnstructured(item)
		if err != nil {
			return err
		}
		err = definition.SetSuspended(true)
		if err != nil {
			return err
		}

		return nil
	}
}

func unsetDowntimeForStream(key string) UnstructuredProcessor {
	return func(def streamapis.Definition) error {
		item := def.ToUnstructured()
		labels := item.GetLabels()

		if labels["arcane.sneaksanddata.com/downtime"] != key {
			return nil // Skip items that don't match the downtime key
		}

		delete(labels, "arcane.sneaksanddata.com/downtime")
		item.SetLabels(labels)

		definition, err := streamapis.FromUnstructured(item)
		if err != nil {
			return err
		}

		err = definition.SetSuspended(false)
		if err != nil {
			return err
		}

		return nil
	}
}

func (s *downtime) runWithQueue(ctx context.Context, streamClass string, filter UnstructuredObjectFilter, process UnstructuredProcessor, printer printers.ResourcePrinter) error {
	rateLimiter := workqueue.DefaultTypedControllerRateLimiter[streamapis.Definition]()
	queue := workqueue.NewTypedRateLimitingQueue[streamapis.Definition](rateLimiter)
	defer queue.ShutDown()
	var wg sync.WaitGroup

	wg.Go(func() {
		s.processObjects(ctx, queue, process, printer)
	})

	err := s.getObjectsList(ctx, streamClass, filter, queue)
	if err != nil {
		return err
	}

	queue.ShutDownWithDrain()
	wg.Wait()
	return nil
}

func (s *downtime) getObjectsList(ctx context.Context, streamClass string, matches UnstructuredObjectFilter, queue Queue) error {
	clientSet, err := s.clientProvider.ProvideClientSet()
	if err != nil {
		return err
	}
	sc, err := clientSet.StreamingV1().StreamClasses("").Get(ctx, streamClass, metav1.GetOptions{})
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

	unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
	if err != nil {
		return err
	}
	err = unstructuredClient.List(ctx, streamList)
	if err != nil {
		return err
	}

	for _, item := range streamList.Items {
		streamDefinition, err := streamapis.FromUnstructured(&item)
		if err != nil {
			logError(&item, "parsing kubernetes object, skipping", err)
			continue // Skip items that can't be parsed as stream definitions
		}
		if !matches(streamDefinition) {
			continue
		}
		queue.Add(streamDefinition)
	}

	return nil
}

func filterByNamePrefix(prefix string) func(streamapis.Definition) bool {
	return func(u streamapis.Definition) bool {
		return strings.HasPrefix(u.ToUnstructured().GetName(), prefix) && !u.Suspended()
	}
}

func filterByDowntimeKey(key string) func(streamapis.Definition) bool {
	return func(u streamapis.Definition) bool {
		return u.ToUnstructured().GetLabels()["arcane.sneaksanddata.com/downtime"] == key
	}
}

func (s *downtime) processObjects(ctx context.Context, queue Queue, process UnstructuredProcessor, printer printers.ResourcePrinter) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			item, shutdown := queue.Get()
			if shutdown {
				return
			}

			err := process(item)
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

			err = unstructuredClient.Update(ctx, item.ToUnstructured())
			if err != nil {
				logError(item.ToUnstructured(), "updating client, will retry later", err)
				queue.AddRateLimited(item)
				continue
			}

			queue.Forget(item)
			queue.Done(item)
			err = printer.PrintObj(item.ToUnstructured(), os.Stdout)
			if err != nil {
				// If we can't print, we still consider the item processed successfully, so we forget it and move on.
				continue
			}
		}
	}
}

func logError(object *unstructured.Unstructured, operation string, cause error) {
	name := FormatName(object)
	_, err := fmt.Fprintf(os.Stderr, "%s Failed %s: %v\n", name, operation, cause)
	if err != nil {
		panic(err)
	}
}
