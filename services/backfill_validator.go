package services

import (
	"context"

	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
)

var _ interfaces.BackfillService = (*backfillValidator)(nil)

// backfillValidator is a service that provides validation for backfill operations.
type backfillValidator struct {
	backfillService interfaces.BackfillService
	clientProvider  interfaces.ClientProvider
}

// NewBackfillValidatorService creates an instance of backfillValidator, which provides validation for backfill operations.
func NewBackfillValidatorService(clientProvider interfaces.ClientProvider) interfaces.BackfillService {
	return &backfillValidator{
		backfillService: &backfill{
			clientProvider: clientProvider,
		},
		clientProvider: clientProvider,
	}
}

func (b *backfillValidator) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {

	return b.backfillService.Backfill(ctx, parameters)
}
