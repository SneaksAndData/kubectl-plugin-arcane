package models

import (
	"github.com/spf13/cobra"
)

// DowntimeDetailsParameters represents the parameters required to perform a list operation for active downtimes.
type DowntimeDetailsParameters struct {
	StreamClass string // The optional stream class filter
}

// NewDowntimeDetailsParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeDetailsParameters(cmd *cobra.Command) (*DowntimeDetailsParameters, error) { // coverage-ignore (tested in integration tests)
	streamClass, err := cmd.Flags().GetString("stream-class")
	if err != nil {
		return nil, err
	}
	return &DowntimeDetailsParameters{StreamClass: streamClass}, nil
}
