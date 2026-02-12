package services

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os/exec"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"strings"
	"testing"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_CreateFailedStream(t *testing.T) {
	name := createTestStreamDefinition(t, false)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)

	streamService := NewStreamService(clientSet)
	err := streamService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.NoError(t, err)
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

func createTestStreamDefinition(t *testing.T, shouldFail bool) string {
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
			Suspended:   false,
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
			RunDuration: "15s",
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
