package models

import (
	"github.com/spf13/cobra"
)

// DowntimeStopParameters represents the parameters required to perform a stop operation for a stream.
type DowntimeStopParameters struct {
	StreamClass string // The class of the stream to stop.
	DowntimeKey string // The unique identifier of the downtime to declare.
}

// NewDowntimeStopParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeStopParameters(_ *cobra.Command, args []string) (*DowntimeStopParameters, error) {
	return &DowntimeStopParameters{
		StreamClass: args[0],
		DowntimeKey: args[1],
	}, nil
}
