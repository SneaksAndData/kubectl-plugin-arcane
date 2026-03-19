package services

import (
	"fmt"
	"strings"
	"testing"
	"time"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDowntime_DeclareDowntime(t *testing.T) {
	// Arrange
	pattern := "declare-downtime-test-"

	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Spec.RunDuration = "5s"
		def.Spec.Suspended = true
		def.Spec.ShouldFail = false
		def.GenerateName = pattern
	})
	require.NotEmpty(t, name)

	err := waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)

	err = WakeUp(t, name)
	require.NoError(t, err)

	downtimeService := createDowntimeService(t)

	// Act
	err = downtimeService.DeclareDowntime(t.Context(), &models.DowntimeDeclareParameters{
		StreamClass: "arcane-stream-mock",
		DowntimeKey: "maintenance-window-1",
		Prefix:      pattern,
	})
	require.NoError(t, err)

	// Assert
	s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)

	err = waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)
	require.Contains(t, s.Labels, interfaces.DowntimeLabelKey)
	require.Contains(t, s.Labels, interfaces.DowntimeBeginLabelKey)
}

func TestDowntime_StopDowntime(t *testing.T) {
	// Arrange
	pattern := "stop-downtime-test-"

	name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Labels = map[string]string{
			interfaces.DowntimeLabelKey:      "maintenance-window-1",
			interfaces.DowntimeBeginLabelKey: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}
		def.Spec.RunDuration = "5s"
		def.Spec.Suspended = true
		def.Spec.ShouldFail = false
		def.GenerateName = pattern
	})
	require.NotEmpty(t, name)

	err := waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)

	downtimeService := createDowntimeService(t)

	// Act
	err = downtimeService.StopDowntime(t.Context(), &models.DowntimeStopParameters{
		StreamClass: "arcane-stream-mock",
		DowntimeKey: "maintenance-window-1",
	})
	require.NoError(t, err)

	// Assert
	err = waitForPhase(t, name, streamapis.Suspended)
	require.NoError(t, err)

	s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
	require.NoError(t, err)
	require.NotContains(t, s.Annotations, interfaces.DowntimeLabelKey)
	require.NotContains(t, s.Annotations, interfaces.DowntimeBeginLabelKey)
	require.False(t, s.Spec.Suspended)
}

func TestDowntime_List_NoFilter(t *testing.T) {
	// Arrange
	const streamCount = 3
	pattern := "list-downtime-test-"

	for i := range streamCount {
		name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
			def.Labels = map[string]string{
				interfaces.DowntimeLabelKey:      fmt.Sprintf("maintenance-window-%d", i),
				interfaces.DowntimeBeginLabelKey: fmt.Sprintf("%d", time.Now().UnixMilli()),
			}
			def.Spec.Suspended = true
			def.GenerateName = pattern
		})
		require.NotEmpty(t, name)
		err := waitForPhase(t, name, streamapis.Suspended)
		require.NoError(t, err)
	}

	downtimeService := createDowntimeService(t)
	dts, err := downtimeService.GetSummary(t.Context(), &models.DowntimeSummaryParameters{
		StreamClass: "",
	})

	require.NoError(t, err)
	require.NotEmpty(t, dts)

	for key, count := range dts.CountsRaw() {
		if strings.HasPrefix(key, "maintenance-window-") {
			require.GreaterOrEqual(t, count, 1)
		}
	}
}

func TestDowntime_Details_NoFilter(t *testing.T) {
	// Arrange
	const streamCount = 3
	pattern := "details-downtime-test-"

	for i := range streamCount {
		name := helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
			def.Labels = map[string]string{
				interfaces.DowntimeLabelKey:      fmt.Sprintf("details-maintenance-window-%d", i),
				interfaces.DowntimeBeginLabelKey: fmt.Sprintf("%d", time.Now().UnixMilli()),
			}
			def.Spec.Suspended = true
			def.GenerateName = pattern
		})
		require.NotEmpty(t, name)
		err := waitForPhase(t, name, streamapis.Suspended)
		require.NoError(t, err)
	}

	downtimeService := createDowntimeService(t)
	dts, err := downtimeService.GetSummary(t.Context(), &models.DowntimeSummaryParameters{
		StreamClass: "",
	})

	require.NoError(t, err)
	require.NotEmpty(t, dts)

	for key, count := range dts.DetailsRaw() {
		if strings.HasPrefix(key, "details-maintenance-window-") {
			require.GreaterOrEqual(t, len(count), 1)
		}
	}
}

func createDowntimeService(t *testing.T) cmdinterfaces.DowntimeService {
	streamingClientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	clientProvider := NewFakeClientProvider(streamingClientSet, c)
	downtimeService := NewDowntimeService(clientProvider, NewDowntimeProcessorFactory(NewUnstructuredReader(clientProvider)))

	return downtimeService
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

	return waitForPhase(t, name, streamapis.Running)
}
