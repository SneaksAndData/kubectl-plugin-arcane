package services

import (
	"context"
	"errors"
	"flag"
	"fmt"
	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os/exec"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"strings"
	"sync"
	"testing"
	"time"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Backfill(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s")
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	streamService := NewStreamService(clientSet, nil)
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
	name := createTestStreamDefinition(t, false, "5s")
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	streamService := NewStreamService(clientSet, nil)
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
	name := createTestStreamDefinition(t, false, "30s")
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel() // Ensure context is cleaned up even if test fails

	streamService := NewStreamService(clientSet, nil)
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

var (
	kubeconfigCmd string
	kubeConfig    *rest.Config
	clientSet     *mockversionedv1.Clientset
)

// TestMain initializes the Kubernetes client and creates a suspended StreamDefinition for testing
func TestMain(m *testing.M) {

	flag.StringVar(&kubeconfigCmd, "kubeconfig-cmd", "/opt/homebrew/bin/kind get kubeconfig", "Command to execute that outputs kubeconfig YAML content")
	flag.Parse()

	// Initialize logger to avoid controller-runtime warnings
	klog.InitFlags(nil)
	logger := klog.Background()
	controllerruntime.SetLogger(logger)

	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
		return
	}

	var err error
	kubeConfig, err = readKubeconfig()
	if err != nil {
		panic(fmt.Errorf("error reading kubeconfig: %w", err))
	}

	clientSet, err = mockversionedv1.NewForConfig(kubeConfig)
	if err != nil {
		panic(fmt.Errorf("error creating kubernetes clientSet: %w", err))
	}

	m.Run()
}

func createTestStreamDefinition(t *testing.T, shouldFail bool, runDuration string) string {
	testStream := &mockv1.TestStreamDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "streaming.sneaksanddata.com/v1",
			Kind:       "TestStreamDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-stream-",
			Namespace:    "default",
		},
		Spec: mockv1.TestsStreamDefinitionSpec{
			Source:      "mock-source",
			Destination: "mock-destination",
			Suspended:   true,
			ShouldFail:  shouldFail,
			JobTemplateRef: corev1.ObjectReference{
				APIVersion: "streaming.sneaksanddata.com/v1",
				Kind:       "StreamingJobTemplate",
				Name:       "arcane-stream-mock",
				Namespace:  "default",
			},
			BackfillJobTemplateRef: corev1.ObjectReference{
				APIVersion: "streaming.sneaksanddata.com/v1",
				Kind:       "StreamingJobTemplate",
				Name:       "arcane-stream-mock",
				Namespace:  "default",
			},
			RunDuration: runDuration,
			TestSecretRef: &corev1.LocalObjectReference{
				Name: "test-secret",
			},
		},
	}

	stream, err := clientSet.
		StreamingV1().
		TestStreamDefinitions(testStream.Namespace).
		Create(t.Context(), testStream, metav1.CreateOptions{})
	require.NoError(t, err)

	return stream.Name
}

func readKubeconfig() (*rest.Config, error) {

	// Parse and execute the command
	cmdParts := strings.Fields(kubeconfigCmd)
	if len(cmdParts) == 0 {
		return nil, errors.New("kubeconfig-cmd cannot be empty")
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("error executing command: %w\nStderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("error executing command: %w", err)
	}

	// Load the kubeconfig from bytes and convert to rest.Config
	clientConfig, err := clientcmd.NewClientConfigFromBytes(output)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %w", err)
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error converting to rest.Config: %w", err)
	}

	return restConfig, nil
}

func findBackfillRequestByName(ctx context.Context, namespace string, name string) (*v1.BackfillRequest, error) {
	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	backfillList, err := clientSet.StreamingV1().BackfillRequests(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing backfill requests: %w", err)
	}
	for _, backfill := range backfillList.Items {
		if backfill.Spec.StreamId == name {
			return &backfill, nil
		}
	}
	return nil, fmt.Errorf("backfill request for stream %s not found in namespace %s", name, namespace)
}
