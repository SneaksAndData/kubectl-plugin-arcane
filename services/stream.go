package services

import (
	"context"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"k8s.io/client-go/rest"
)

// Ensure stream implements interfaces.StreamService
var _ interfaces.StreamService = (*stream)(nil)

// stream is a service that provides stream operations.
type stream struct {
	restConfig *rest.Config
}

// Backfill is a method that allows users to run a stream backfill operation, use the <key> parameter to identify the stream to backfill
func (s *stream) Backfill(ctx context.Context, parameters *models.BackfillParameters) error {
	panic("not implemented")
}

// NewStreamService creates a new instance of the stream, which provides stream operations.
func NewStreamService(config rest.Config) interfaces.StreamService {
	return &stream{
		restConfig: &config,
	}
}

// Start is a method that allows users to start a stream, use the <key> parameter to identify the stream to start
func (s *stream) Start(ctx context.Context, parameters *models.StartParameters) error {
	panic("not implemented")
}

// Stop is a method that allows users to stop a stream, use the <key> parameter to identify the stream to stop
func (s *stream) Stop(ctx context.Context, parameters *models.StopParameters) error {
	panic("not implemented")
}
