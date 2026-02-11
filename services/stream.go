package services

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/spf13/cobra"
)

// Ensure stream implements interfaces.StreamService
var _ interfaces.StreamService = (*stream)(nil)

// stream is a service that provides stream operations.
type stream struct {
}

// Backfill is a method that allows users to run a stream backfill operation, use the <key> parameter to identify the stream to backfill
func (s *stream) Backfill(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}

// NewStreamService creates a new instance of the stream, which provides stream operations.
func NewStreamService() interfaces.StreamService {
	return &stream{}
}

// Execute is a method that allows users to run a stream command, use the <key> parameter to identify the stream to execute
func (s *stream) Execute(cmd *cobra.Command, args []string) {
	panic("not implemented")
}

// Start is a method that allows users to start a stream, use the <key> parameter to identify the stream to start
func (s *stream) Start(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}

// Stop is a method that allows users to stop a stream, use the <key> parameter to identify the stream to stop
func (s *stream) Stop(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}
