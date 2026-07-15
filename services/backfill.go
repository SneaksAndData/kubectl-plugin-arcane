package services

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ interfaces.BackfillService = (*backfill)(nil)

// backfill is a service that provides backfill operations.
type backfill struct {
	clientProvider interfaces.ClientProvider
}

// NewBackfillService creates a new instance of the backfill, which provides backfill operations.
func NewBackfillService(clientProvider interfaces.ClientProvider) interfaces.BackfillService {
	return &backfill{
		clientProvider: clientProvider,
	}
}

// Backfill is a method that allows users to run a stream backfill operation
func (b *backfill) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {
	clientSet, err := b.clientProvider.ProvideClientSet()
	if err != nil {
		return fmt.Errorf("error providing client set: %w", err)
	}

	bfr, err := b.getBackfillRequest(ctx, clientSet, parameters.Namespace, parameters.StreamId)
	if err != nil {
		return fmt.Errorf("error checking for existence of an backfill request: %w", err)
	}
	if bfr == nil {
		bfr, err = clientSet.
			StreamingV1().
			BackfillRequests(parameters.Namespace).
			Create(ctx, parameters.ToBackfillRequest(), metav1.CreateOptions{
				FieldManager:    fieldManager,
				FieldValidation: "Strict",
			})
		if err != nil {
			return fmt.Errorf("error creating backfill request: %w", err)
		}

		if !parameters.Wait {
			return logging.Printer("created").PrintObj(bfr, os.Stdout)
		}
	} else {
		return logging.Printer("already exists").PrintObj(bfr, os.Stdout)
	}

	err = logging.Printer("started").PrintObj(bfr, os.Stdout)
	if err != nil {
		return err
	}
	return wait.PollUntilContextCancel(ctx, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		watch, err := clientSet.StreamingV1().BackfillRequests(parameters.Namespace).Watch(ctx, metav1.ListOptions{
			FieldSelector:   "metadata.name=" + bfr.Name,
			Watch:           true,
			ResourceVersion: bfr.ResourceVersion,
		})
		if err != nil {
			return false, fmt.Errorf("error watching backfill request: %w", err)
		}
		defer watch.Stop()

		for {
			select {
			case event, ok := <-watch.ResultChan():
				if !ok {
					logging.LogError(bfr, "watching backfill request, retrying", fmt.Errorf("watch channel closed"))
					return false, nil // watch channel closed, retry
				}
				bfr, ok := event.Object.(*v1.BackfillRequest)
				if !ok {
					return false, fmt.Errorf("unexpected object type: %T", event.Object)
				}

				if bfr.Spec.Completed {
					return true, logging.Printer("completed").PrintObj(bfr, os.Stdout)
				}

			case <-ctx.Done():
				return true, ctx.Err()
			}
		}
	})
}

func (b *backfill) getBackfillRequest(ctx context.Context, clientSet *versioned.Clientset, namespace string, id string) (*v1.BackfillRequest, error) {
	list, err := clientSet.
		StreamingV1().
		BackfillRequests(namespace).
		List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.completed=false,spec.streamId=%s", id),
		})

	if err != nil {
		return nil, fmt.Errorf("error listing backfill requests: %w", err)
	}

	if len(list.Items) > 1 {
		return nil, fmt.Errorf("multiple active backfill requests found for stream %s in namespace %s", namespace, namespace)
	}

	if len(list.Items) == 0 {
		return nil, nil
	}

	return &list.Items[0], nil
}
