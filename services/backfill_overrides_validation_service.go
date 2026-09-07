package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/data_structures"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ interfaces.BackfillService = (*backfillOverridesValidationService)(nil)

// backfillOverridesValidationService is a service that provides validation for backfill operations.
type backfillOverridesValidationService struct {
	backfillService interfaces.BackfillService
	clientProvider  interfaces.ClientProvider
	pt              *data_structures.PrefixTree
}

// newBackfillOverridesValidationService creates an instance of backfillOverridesValidationService, which provides validation for backfill operations.
func newBackfillOverridesValidationService(clientProvider interfaces.ClientProvider) interfaces.BackfillService {
	return &backfillOverridesValidationService{
		backfillService: newBackfillService(clientProvider),
		clientProvider:  clientProvider,
		pt:              data_structures.NewPrefixTree(),
	}
}

func (b *backfillOverridesValidationService) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {
	clientSet, err := b.clientProvider.ProvideClientSet()
	if err != nil {
		return fmt.Errorf("backfillOverridesValidationService: error creating clientSet: %w", err)
	}

	streamClass, err := clientSet.StreamingV1().StreamClasses("").Get(ctx, parameters.StreamClass, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("backfillOverridesValidationService: error getting stream class: %w", err)
	}

	err = b.pt.InsertAll(streamClass.Spec.OverridableFields)
	if err != nil {
		return fmt.Errorf("backfillOverridesValidationService: error inserting overrides into the trie: %w", err)
	}

	if parameters.Overrides != nil {
		for _, override := range *parameters.Overrides {
			key := strings.SplitN(override, "=", 2)[0]
			hasPrefix, err := b.pt.HasPrefix(key)
			if err != nil {
				return fmt.Errorf("backfillOverridesValidationService: error checking prefix: %w", err)
			}
			if !hasPrefix {
				return fmt.Errorf("backfillOverridesValidationService: parameter is not overridable: %s", key)
			}
		}
	}

	return b.backfillService.Backfill(ctx, parameters)
}
