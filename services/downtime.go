package services

import (
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/spf13/cobra"
)

// Ensure downtime implements interfaces.DowntimeService
var _ interfaces.DowntimeService = (*downtime)(nil)

// downtime is a service that provides downtime operations.
type downtime struct {
}

// NewDowntimeService creates a new instance of the downtime, which provides downtime operations.
func NewDowntimeService() interfaces.DowntimeService {
	return &downtime{}
}

// DeclareDowntime is a method that allows users to declare downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to pause
func (s *downtime) DeclareDowntime(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}

// StopDowntime is a method that allows users to stop downtime for a stream or a list of streams, use the <key> parameter to identify the stream(s) to resume
func (s *downtime) StopDowntime(cmd *cobra.Command, args []string) error {
	panic("not implemented")
}
