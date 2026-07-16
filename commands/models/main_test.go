package models

import (
	"flag"
	"fmt"
	"testing"

	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var (
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

	_, err = mockversionedv1.NewForConfig(kubeConfig)
	if err != nil {
		panic(fmt.Errorf("error creating kubernetes clientSet: %w", err))
	}

	m.Run()
}
