package models

import (
	"github.com/spf13/cobra"
)

// DowntimeListParameters represents the parameters required to perform a list operation for active downtimes.
type DowntimeListParameters struct {
	StreamClass *string // The class of the stream to stop.
}

// NewDowntimeListParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeListParameters(cmd *cobra.Command) (*DowntimeListParameters, error) { // coverage-ignore (tested in integration tests)
	streamClass, err := cmd.Flags().GetString("stream-class")
	if err != nil {
		return nil, err
	}
	if streamClass == "" {
		return &DowntimeListParameters{StreamClass: nil}, nil
	}
	return &DowntimeListParameters{StreamClass: &streamClass}, nil
}
