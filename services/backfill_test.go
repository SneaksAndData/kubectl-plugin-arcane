package services

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Backfill(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	backfillService := NewBackfillService(NewFakeClientProvider(clientSet, nil))
	err := backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.NoError(t, err)
	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.NoError(t, err)
	require.False(t, bfr.Spec.Completed)
}

func Test_Backfill_Wait(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	backfillService := NewBackfillService(NewFakeClientProvider(clientSet, nil))
	err := backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        true,
	})
	require.NoError(t, err)
	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.NoError(t, err)
	require.True(t, bfr.Spec.Completed)
}

func Test_Backfill_Duplicate(t *testing.T) {
	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Spec.RunDuration = "30m"
		def.Spec.Suspended = false
	})
	require.NotEmpty(t, name)

	err := waitForPhase(t, name, streamapis.Backfilling)
	require.NoError(t, err)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	backfillService := NewBackfillService(NewFakeClientProvider(clientSet, nil))

	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.NoError(t, err)

	backfillList, err := clientSet.StreamingV1().BackfillRequests("default").List(t.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.streamId=%s", name),
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(backfillList.Items))
}

func Test_Backfill_CompletedDuplicate(t *testing.T) {
	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Spec.RunDuration = "3s"
		def.Spec.Suspended = false
	})
	require.NotEmpty(t, name)

	err := waitForPhase(t, name, streamapis.Running)
	require.NoError(t, err)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	backfillService := NewBackfillService(NewFakeClientProvider(clientSet, nil))

	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.NoError(t, err)

	err = waitForPhase(t, name, streamapis.Running)
	require.NoError(t, err)

	backfillList, err := clientSet.StreamingV1().BackfillRequests("default").List(t.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.streamId=%s", name),
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(backfillList.Items))
}

func Test_Backfill_Cancelled(t *testing.T) {
	name := createTestStreamDefinition(t, false, "30s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel() // Ensure context is cleaned up even if test fails

	backfillService := NewBackfillService(NewFakeClientProvider(clientSet, nil))
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // Ensure Done() is called even if Backfill panics
		err = backfillService.Backfill(ctx, &models.BackfillParameters{
			Namespace:   "default",
			StreamId:    name,
			StreamClass: "arcane-stream-mock",
			Wait:        true,
		})
	}()

	// Cancel the context to simulate cancellation during backfill
	time.Sleep(5 * time.Second)
	cancel()

	wg.Wait()

	// Expect context.Canceled error
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	// Verify that backfill request was created but not completed
	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.NoError(t, err)
	require.False(t, bfr.Spec.Completed, "backfill should not be completed when context is cancelled")
}
