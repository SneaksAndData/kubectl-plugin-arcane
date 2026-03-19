package models

import (
	"github.com/spf13/cobra"
)

// DowntimeSummaryParameters represents the parameters required to perform a list operation for active downtimes.
type DowntimeSummaryParameters struct {
	StreamClass string // The optional stream class filter
}

// NewDowntimeListParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeListParameters(cmd *cobra.Command) (*DowntimeSummaryParameters, error) { // coverage-ignore (tested in integration tests)
	streamClass, err := cmd.Flags().GetString("stream-class")
	if err != nil {
		return nil, err
	}
	return &DowntimeSummaryParameters{StreamClass: streamClass}, nil
}
