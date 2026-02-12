package models

import (
	"fmt"
	"github.com/spf13/cobra"
)

// StartParameters represents the parameters required to perform a stop operation for a stream.
type StartParameters struct {
	StreamClass string // The class of the stream to stop.
	StreamID    string // The unique identifier of the stream to stop.
}

// NewStartParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewStartParameters(_ *cobra.Command, args []string) (*StartParameters, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid arguments for stop parameters")
	}

	return &StartParameters{StreamClass: args[0], StreamID: args[1]}, nil
}
