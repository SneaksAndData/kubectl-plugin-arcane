package services

import (
	"context"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
	"testing"
	"time"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Backfill(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	streamService := NewStreamService(NewFakeClientProvider(clientSet, nil))
	err := streamService.Backfill(t.Context(), &models.BackfillParameters{
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

	streamService := NewStreamService(NewFakeClientProvider(clientSet, nil))
	err := streamService.Backfill(t.Context(), &models.BackfillParameters{
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

func Test_Backfill_Cancelled(t *testing.T) {
	name := createTestStreamDefinition(t, false, "30s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel() // Ensure context is cleaned up even if test fails

	streamService := NewStreamService(NewFakeClientProvider(clientSet, nil))
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // Ensure Done() is called even if Backfill panics
		err = streamService.Backfill(ctx, &models.BackfillParameters{
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

func Test_StreamStarted(t *testing.T) {
	name := createTestStreamDefinition(t, false, "15s", true)
	require.NotEmpty(t, name)

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	streamService := NewStreamService(NewFakeClientProvider(streamingClientSet, c))
	err = streamService.Start(t.Context(), &models.StartParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
	})
	require.NoError(t, err)

	stream, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.False(t, stream.Spec.Suspended)
}

func Test_StreamStopped(t *testing.T) {
	name := createTestStreamDefinition(t, false, "15s", false)
	require.NotEmpty(t, name)
	err := wait.PollUntilContextCancel(t.Context(), 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return s.Status.Phase == string(streamapis.Running), nil
	})

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	streamService := NewStreamService(NewFakeClientProvider(streamingClientSet, c))
	err = streamService.Stop(t.Context(), &models.StopParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
	})
	require.NoError(t, err)

	err = wait.PollUntilContextCancel(t.Context(), 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
		return s.Spec.Suspended, err
	})

	require.NoError(t, err)
	stream, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, stream.Spec.Suspended)
}
