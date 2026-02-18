package tests

import (
	"context"
	"flag"
	"fmt"
	mockv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v1"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test_Start(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-start-"
		},
		"kubectl arcane stream start arcane-stream-mock %s --namespace integration-tests",
	)
}

func Test_Stop(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = false
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-stop-"
		},
		"kubectl arcane stream stop arcane-stream-mock %s --namespace integration-tests",
	)
}

func Test_Backfill(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-backfill-"
		},
		"kubectl arcane stream backfill arcane-stream-mock %s --namespace integration-tests",
	)
}

func Test_Backfill_Wait(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-backfill-wait-"
		},
		"kubectl arcane stream backfill arcane-stream-mock %s --wait --namespace integration-tests",
	)
}

func Test_DowntimeDeclare(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = false
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-declare-"
		},
		"kubectl arcane downtime declare arcane-stream-mock %s downtime-window-1 --namespace integration-tests",
	)
}

func Test_DowntimeStop(t *testing.T) {
	runIntegrationTest(t,
		func(def *mockv1.TestStreamDefinition) {
			def.Namespace = "integration-tests"
			def.Labels = map[string]string{
				"arcane.sneaksanddata.com/downtime": "maintenance-window-1",
			}
			def.Spec.RunDuration = "5s"
			def.Spec.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-declare-"
		},
		"kubectl arcane downtime stop arcane-stream-mock downtime-window-1 --namespace integration-tests",
	)
}

var clientSet *mockversionedv1.Clientset

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode.")
		os.Exit(0)
	}
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Add bin directory to PATH
	binDir := filepath.Join(currentDir, "bin")
	currentPath := os.Getenv("PATH")
	newPath := binDir + string(filepath.ListSeparator) + currentPath
	err = os.Setenv("PATH", newPath)
	if err != nil {
		panic(err)
	}

	var kubeconfigCmd string
	flag.StringVar(&kubeconfigCmd, "kubeconfig-cmd", "/opt/homebrew/bin/kind get kubeconfig", "Command to execute that outputs kubeconfig YAML content")
	flag.Parse()

	kubeconfigContent, err := helpers.GetKubeconfigString(kubeconfigCmd)
	if err != nil {
		panic(fmt.Errorf("error reading kubeconfig: %w", err))
	}

	// Write kubeconfig to bin directory
	kubeconfigPath := filepath.Join(binDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, kubeconfigContent, 0600)
	if err != nil {
		panic(fmt.Errorf("error writing kubeconfig: %w", err))
	}

	kubeConfig, err := helpers.DeserializeKubeconfig(kubeconfigContent)
	if err != nil {
		panic(fmt.Errorf("error deserializing kubeconfig: %w", err))
	}

	clientSet, err = mockversionedv1.NewForConfig(kubeConfig)
	require.NoError(nil, err, "error creating kubernetes clientProvider")

	// Set KUBECONFIG environment variable
	err = os.Setenv("KUBECONFIG", kubeconfigPath)
	if err != nil {
		panic(err)
	}

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func runCommand(ctx context.Context, args string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", args)
	return cmd.CombinedOutput()
}

func runIntegrationTest(t *testing.T, setup func(def *mockv1.TestStreamDefinition), commandTemplate string) {
	name := helpers.NewTestStream(t, clientSet, setup)

	require.NotEmpty(t, name)
	var command string
	if strings.Contains(commandTemplate, "%s") {
		command = fmt.Sprintf(commandTemplate, name)
	} else {
		command = commandTemplate
	}
	fmt.Println(command)
	output, err := runCommand(t.Context(), command)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
	}
	t.Logf("Command output:\n%s", string(output))
}
