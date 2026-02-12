package models

import (
	"fmt"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackfillParameters represents the parameters required to perform a backfill operation for a stream.
type BackfillParameters struct {
	StreamClass string   // The class of the stream to backfill.
	StreamId    string   // The unique identifier of the stream to backfill.
	Wait        bool     // Whether to wait for the backfill operation to complete before returning.
	Namespace   string   // The Kubernetes Namespace where the stream is located.
	DryRun      []string // The dry-run strategy to apply when creating the backfill request (e.g., "All", "Client", "Server").
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

	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return nil, err
	}

	dryRun, err := cmd.Flags().GetBool("dryRun")
	if err != nil {
		return nil, err
	}

	var dryRunStrategy []string
	if dryRun {
		dryRunStrategy = []string{"All"}
	}

	bfr := &BackfillParameters{
		StreamClass: args[0],
		StreamId:    args[1],
		Wait:        wait,
		Namespace:   namespace,
		DryRun:      dryRunStrategy,
	}

	return bfr, nil
}

func (p BackfillParameters) ToBackfillRequest() *v1.BackfillRequest {
	return &v1.BackfillRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    p.Namespace,
			GenerateName: fmt.Sprintf("%s-manual-", p.StreamId),
		},
		Spec: v1.BackfillRequestSpec{
			StreamClass: p.StreamClass,
			StreamId:    p.StreamId,
		},
	}
}
