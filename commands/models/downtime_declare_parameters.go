package models

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// DowntimeDeclareParameters represents the parameters required to perform a stop operation for a stream.
type DowntimeDeclareParameters struct {
	StreamClass string // The class of the stream to stop.
	Prefix      string // The prefix of the stream to stop.
	DowntimeKey string // The unique identifier of the downtime to declare.
	Namespace   string // The namespace of the stream to stop.
}

// NewDowntimeDeclareParameters creates a new instance of StopParameters based on the provided command and arguments.
func NewDowntimeDeclareParameters(_ *cobra.Command, args []string, configFlags *genericclioptions.ConfigFlags) (*DowntimeDeclareParameters, error) { // coverage-ignore (tested in integration tests)
	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}
	return &DowntimeDeclareParameters{
		StreamClass: args[0],
		Prefix:      args[1],
		DowntimeKey: args[2],
		Namespace:   namespace,
	}, nil
}
