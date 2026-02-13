package models

import (
	"fmt"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// BackfillParameters represents the parameters required to perform a backfill operation for a stream.
type BackfillParameters struct {
	StreamClass string // The class of the stream to backfill.
	StreamId    string // The unique identifier of the stream to backfill.
	Wait        bool   // Whether to wait for the backfill operation to complete before returning.
	Namespace   string // The namespace in which the stream is located. If empty, the default namespace will be used.
}

// NewBackfillParameters creates a new instance of BackfillParameters based on the provided command and arguments.
func NewBackfillParameters(cmd *cobra.Command, args []string, configFlags *genericclioptions.ConfigFlags) (*BackfillParameters, error) {

	wait, err := cmd.Flags().GetBool("wait")
	if err != nil {
		return nil, err
	}

	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	bfr := &BackfillParameters{
		StreamClass: args[0],
		StreamId:    args[1],
		Wait:        wait,
		Namespace:   namespace,
	}

	return bfr, nil
}

func (p BackfillParameters) ToBackfillRequest() *v1.BackfillRequest {
	return &v1.BackfillRequest{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-manual-", p.StreamId),
		},
		Spec: v1.BackfillRequestSpec{
			StreamClass: p.StreamClass,
			StreamId:    p.StreamId,
		},
	}
}
