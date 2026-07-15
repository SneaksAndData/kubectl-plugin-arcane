package services

import (
	"testing"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_StreamStarted(t *testing.T) {
	name := createTestStreamDefinition(t, false, "15s", true)
	require.NotEmpty(t, name)

	err := waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)
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

func Test_StreamStarted_Error(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s", false)
	require.NotEmpty(t, name)
	err := waitForPhase(t, name, streamapis.Running)
	require.NoError(t, err)

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	streamService := NewStreamService(NewFakeClientProvider(streamingClientSet, c))
	err = streamService.Start(t.Context(), &models.StartParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Stream already has desired phase Running")
}

func Test_StreamStopped(t *testing.T) {
	name := createTestStreamDefinition(t, false, "15s", false)
	require.NotEmpty(t, name)
	err := waitForPhase(t, name, streamapis.Running)
	require.NoError(t, err)

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

	stream, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.True(t, stream.Spec.Suspended)
}

func Test_StreamStopped_Error(t *testing.T) {
	name := createTestStreamDefinition(t, false, "15s", true)
	require.NotEmpty(t, name)
	err := waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	streamService := NewStreamService(NewFakeClientProvider(streamingClientSet, c))
	err = streamService.Stop(t.Context(), &models.StopParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Stream already has desired phase Suspended")
}
