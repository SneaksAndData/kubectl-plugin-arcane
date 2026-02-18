package services

import (
	"context"
	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

func TestDowntime_DeclareDowntime(t *testing.T) {
	pattern := "declare-downtime-test-"

	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Spec.RunDuration = "5s"
		def.Spec.Suspended = true
		def.Spec.ShouldFail = false
		def.GenerateName = pattern
	})
	require.NotEmpty(t, name)

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	downtimeService := NewDowntimeService(streamingClientSet, c)

	err = WakeUp(t, name)
	require.NoError(t, err)

	err = downtimeService.DeclareDowntime(t.Context(), &models.DowntimeDeclareParameters{
		StreamClass: "arcane-stream-mock",
		DowntimeKey: "maintenance-window-1",
		Prefix:      pattern,
	})
	require.NoError(t, err)

	s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.Contains(t, s.Labels, "arcane.sneaksanddata.com/downtime")
	require.True(t, s.Spec.Suspended)
}

func TestDowntime_StopDowntime(t *testing.T) {
	pattern := "stop-downtime-test-"

	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Labels = map[string]string{
			"arcane.sneaksanddata.com/downtime": "maintenance-window-1",
		}
		def.Spec.RunDuration = "5s"
		def.Spec.Suspended = true
		def.Spec.ShouldFail = false
		def.GenerateName = pattern
	})
	require.NotEmpty(t, name)

	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	downtimeService := NewDowntimeService(streamingClientSet, c)

	err = downtimeService.StopDowntime(t.Context(), &models.DowntimeStopParameters{
		StreamClass: "arcane-stream-mock",
		DowntimeKey: "maintenance-window-1",
	})
	require.NoError(t, err)

	s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.NotContains(t, s.Annotations, "arcane.sneaksanddata.com/downtime")
	require.False(t, s.Spec.Suspended)
}

func WakeUp(t *testing.T, name string) error {
	// First, patch the stream to unsuspend it
	patchData := []byte(`{"spec":{"suspended":false}}`)
	_, err := clientSet.StreamingV1().TestStreamDefinitions("default").Patch(
		t.Context(),
		name,
		types.MergePatchType,
		patchData,
		metav1.PatchOptions{},
	)
	if err != nil {
		return err
	}

	// Then wait for it to be running
	return wait.PollUntilContextCancel(t.Context(), 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		stream, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return stream.Status.Phase == "Running", nil
	})
}
