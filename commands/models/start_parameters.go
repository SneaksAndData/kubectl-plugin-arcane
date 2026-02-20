package models

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// StartParameters represents the parameters required to perform a stop operation for a stream.
type StartParameters struct {
	StreamClass string // The class of the stream to stop.
	StreamId    string // The unique identifier of the stream to stop.
	Namespace   string // The unique identifier of the stream to stop.
}

// NewStartParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewStartParameters(_ *cobra.Command, args []string, configFlags *genericclioptions.ConfigFlags) (*StartParameters, error) { // coverage-ignore (tested in integration tests)
	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	return &StartParameters{StreamClass: args[0], StreamId: args[1], Namespace: namespace}, nil
}
