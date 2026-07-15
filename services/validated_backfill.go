package services

import (
	"context"
	"fmt"

	"github.com/SneaksAndData/arcane-operator/services/controllers/contracts"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ interfaces.BackfillService = (*validatedBackfill)(nil)

// validatedBackfill is a service that provides validation for backfill operations.
type validatedBackfill struct {
	backfillService interfaces.BackfillService
	clientProvider  interfaces.ClientProvider
}

// NewValidatedBackfillService creates an instance of validatedBackfill, which provides validation for backfill operations.
func NewValidatedBackfillService(clientProvider interfaces.ClientProvider) interfaces.BackfillService {
	return &validatedBackfill{
		backfillService: newBackfillService(clientProvider),
		clientProvider:  clientProvider,
	}
}

func (b *validatedBackfill) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {
	clientSet, err := b.clientProvider.ProvideClientSet()
	if err != nil {
		return fmt.Errorf("validatedBackfill: error creating clientSet: %w", err)
	}

	streamClass, err := clientSet.StreamingV1().StreamClasses("").Get(ctx, parameters.StreamClass, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("validatedBackfill: error getting stream class: %w", err)
	}

	unstructuredClient, err := b.clientProvider.ProvideUnstructuredClient()
	if err != nil {
		return fmt.Errorf("validatedBackfill: error creating unstructured client: %w", err)
	}

	namespacedName := types.NamespacedName{
		Namespace: parameters.Namespace,
		Name:      parameters.StreamId,
	}
	_, err = streamapis.GetStreamForClass(ctx, unstructuredClient, streamClass, namespacedName, contracts.FromUnstructured)
	if err != nil {
		return fmt.Errorf("error fetching stream definition: %w", err)
	}

	return b.backfillService.Backfill(ctx, parameters)
}
