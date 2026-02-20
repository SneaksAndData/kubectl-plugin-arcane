package models

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StopParameters represents the parameters required to perform a stop operation for a stream.
type StopParameters struct {
	StreamClass string // The class of the stream to stop.
	StreamId    string // The unique identifier of the stream to stop.
	Namespace   string // The unique identifier of the stream to stop.
}

// NewStopParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewStopParameters(_ *cobra.Command, args []string, configFlags *genericclioptions.ConfigFlags) (*StopParameters, error) { // coverage-ignore (tested in integration tests)
	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	return &StopParameters{StreamClass: args[0], StreamId: args[1], Namespace: namespace}, nil
}
