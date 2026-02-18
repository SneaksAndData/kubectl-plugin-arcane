package services

import (
	"context"
	"flag"
	"fmt"
	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

var _ interfaces.ClientProvider = (*FakeClientProvider)(nil)

type FakeClientProvider struct {
	clientSet          *versionedv1.Clientset
	unstructuredClient client.Client
}

func (f FakeClientProvider) ProvideClientSet() (*versionedv1.Clientset, error) {
	return f.clientSet, nil
}

func (f FakeClientProvider) ProvideUnstructuredClient() (client.Client, error) {
	return f.unstructuredClient, nil
}

func NewFakeClientProvider(clientSet *versionedv1.Clientset, unstructuredClient client.Client) *FakeClientProvider {
	return &FakeClientProvider{
		clientSet:          clientSet,
		unstructuredClient: unstructuredClient,
	}
}

var (
	clientSet     *mockversionedv1.Clientset
	kubeconfigCmd string
	kubeConfig    *rest.Config
)

// TestMain initializes the Kubernetes client and creates a suspended StreamDefinition for testing
func TestMain(m *testing.M) {

	flag.StringVar(&kubeconfigCmd, "kubeconfig-cmd", "/opt/homebrew/bin/kind get kubeconfig", "Command to execute that outputs kubeconfig YAML content")
	flag.Parse()

	// Initialize logger to avoid controller-runtime warnings
	klog.InitFlags(nil)
	logger := klog.Background()
	controllerruntime.SetLogger(logger)

	var err error
	kubeConfig, err = helpers.ReadKubeconfig(kubeconfigCmd)
	if err != nil {
		panic(fmt.Errorf("error reading kubeconfig: %w", err))
	}

	clientSet, err = mockversionedv1.NewForConfig(kubeConfig)
	if err != nil {
		panic(fmt.Errorf("error creating kubernetes clientProvider: %w", err))
	}

	m.Run()
}

func createTestStreamDefinition(t *testing.T, shouldFail bool, runDuration string, suspended bool) string {
	return helpers.NewTestStream(t, clientSet, func(def *mockv1.TestStreamDefinition) {
		def.Spec.ShouldFail = shouldFail
		def.Spec.RunDuration = runDuration
		def.Spec.Suspended = suspended
	})
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

func waitForPhase(t *testing.T, name string, phase streamapis.Phase) error {
	return wait.PollUntilContextCancel(t.Context(), 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		s, err := clientSet.StreamingV1().TestStreamDefinitions("default").Get(t.Context(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return s.Status.Phase == string(phase), nil
	})
}
