package models

import (
	"fmt"
	"github.com/spf13/cobra"
)

// BackfillParameters represents the parameters required to perform a backfill operation for a stream.
type BackfillParameters struct {
	StreamClass string // The class of the stream to backfill.
	StreamID    string // The unique identifier of the stream to backfill.
	Wait        bool   // Whether to wait for the backfill operation to complete before returning.
}

// NewBackfillParameters creates a new instance of BackfillParameters based on the provided command and arguments.
func NewBackfillParameters(cmd *cobra.Command, args []string) (*BackfillParameters, error) {

	if len(args) != 2 {
		return nil, fmt.Errorf("invalid arguments for backfill parameters")
	}

	wait, err := cmd.Flags().GetBool("wait")
	if err != nil {
		return nil, err
	}

	return &BackfillParameters{StreamClass: args[0], StreamID: args[1], Wait: wait}, nil
}
