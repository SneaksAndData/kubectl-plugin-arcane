package helpers

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v2"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewTestStream(t *testing.T, clientSet *mockversionedv1.Clientset, configure func(v2 *v2.TestStreamDefinitionV2)) string {
	testStream := v2.TestStreamDefinitionV2{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "streaming.sneaksanddata.com/v1",
			Kind:       "TestStreamDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-stream-",
			Namespace:    "default",
		},
		Spec: v2.TestsStreamDefinitionSpec{
			Source:      "mock-source",
			Destination: "mock-destination",
			ExecutionSettings: v2.ExecutionSettings{
				BackfillJobTemplateRef: nil,
				LayoutVersion:          "v2",
				StreamingBackend: v2.StreamingBackend{
					BatchJobBackend: &v2.BatchJobBackend{
						JobTemplateRef: corev1.ObjectReference{
							APIVersion: "streaming.sneaksanddata.com/v1",
							Kind:       "StreamingJobTemplate",
							Name:       "arcane-stream-mock",
							Namespace:  "default",
						},
						BackfillJobTemplateRef: &corev1.ObjectReference{
							APIVersion: "streaming.sneaksanddata.com/v1",
							Kind:       "StreamingJobTemplate",
							Name:       "arcane-stream-mock",
							Namespace:  "default",
						},
					},
				},
			},
			RunDuration: "5s",
			TestSecretRef: &corev1.LocalObjectReference{
				Name: "test-secret",
			},
		},
	}

	configure(&testStream)

	stream, err := clientSet.
		StreamingV2().
		TestStreamDefinitionV2s(testStream.Namespace).
		Create(t.Context(), &testStream, metav1.CreateOptions{})
	require.NoError(t, err)

	return stream.Name
}

func UnsuspendTestStreamDefinition(ctx context.Context, clientSet *mockversionedv1.Clientset, namespace string, name string) error {
	testStreamDefinition, err := clientSet.StreamingV2().TestStreamDefinitionV2s(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error reading test stream definition %s/%s: %w", namespace, name, err)
	}

	testStreamDefinition.Spec.ExecutionSettings.Suspended = false
	_, err = clientSet.StreamingV2().TestStreamDefinitionV2s(namespace).Update(ctx, testStreamDefinition, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating test stream definition %s/%s: %w", namespace, name, err)
	}

	return nil
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
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
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
