package helpers

import (
	"errors"
	"fmt"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os/exec"
	"strings"
	"testing"
)

func NewTestStream(t *testing.T, clientSet *mockversionedv1.Clientset, configure func(*mockv1.TestStreamDefinition)) string {
	testStream := mockv1.TestStreamDefinition{
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
			ShouldFail:  false,
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
			RunDuration: "5s",
			TestSecretRef: &corev1.LocalObjectReference{
				Name: "test-secret",
			},
		},
	}

	configure(&testStream)

	stream, err := clientSet.
		StreamingV1().
		TestStreamDefinitions(testStream.Namespace).
		Create(t.Context(), &testStream, metav1.CreateOptions{})
	require.NoError(t, err)

	return stream.Name
}

func GetKubeconfigString(kubeconfigCmd string) ([]byte, error) {
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

	return output, nil
}

func DeserializeKubeconfig(kubeconfigBytes []byte) (*rest.Config, error) {

	// Load the kubeconfig from bytes and convert to rest.Config
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %w", err)
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error converting to rest.Config: %w", err)
	}

	return restConfig, nil
}

func ReadKubeconfig(kubeconfigCmd string) (*rest.Config, error) {
	output, err := GetKubeconfigString(kubeconfigCmd)
	if err != nil {
		return nil, fmt.Errorf("error getting kubeconfig string: %w", err)
	}

	return DeserializeKubeconfig(output)
}
