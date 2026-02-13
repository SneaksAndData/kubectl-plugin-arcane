package models

import (
	"fmt"
	"github.com/spf13/cobra"
)

// DowntimeDeclareParameters represents the parameters required to perform a stop operation for a stream.
type DowntimeDeclareParameters struct {
	StreamClass string // The class of the stream to stop.
	DowntimeKey string // The unique identifier of the downtime to declare.
	Prefix      string // The prefix of the stream to stop.
	Namespace   string // The unique identifier of the stream to stop.
}

// NewDowntimeDeclareParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeDeclareParameters(_ *cobra.Command, args []string) (*DowntimeDeclareParameters, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid arguments for stop parameters")
	}

	return &DowntimeDeclareParameters{StreamClass: args[0], DowntimeKey: args[1]}, nil
}
