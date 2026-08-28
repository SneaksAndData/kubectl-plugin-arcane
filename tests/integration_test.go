package tests

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v2"
	mockversionedv1 "github.com/SneaksAndData/arcane-stream-mock/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func Test_Start(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-start-"
		},
		nil,
		"kubectl arcane stream start arcane-stream-mock-v2 %s --namespace integration-tests",
	)
}

func Test_Stop(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-stop-"
		},
		nil,
		"kubectl arcane stream stop arcane-stream-mock-v2 %s --namespace integration-tests",
	)
}

func Test_Backfill(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-backfill-"
		},
		nil,
		"kubectl arcane stream backfill arcane-stream-mock-v2 %s --namespace integration-tests",
	)
}

func Test_Backfill_Wait(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-backfill-wait-"
		},
		func(ctx context.Context, clientSet *mockversionedv1.Clientset, name string, namespace string) error {
			return helpers.UnsuspendTestStreamDefinition(ctx, clientSet, name, namespace)
		},
		"kubectl arcane stream backfill arcane-stream-mock-v2 %s --wait --namespace integration-tests",
	)
}

func Test_DowntimeDeclare(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "5s"
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-declare-"
		},
		nil,
		"kubectl arcane downtime declare arcane-stream-mock-v2 %s downtime-window-1 --namespace integration-tests",
	)
}

func Test_DowntimeStop(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Labels = map[string]string{
				interfaces.DowntimeLabelKey: "maintenance-window-1",
			}
			def.Annotations = map[string]string{
				interfaces.DowntimeBeginAnnotationKey: time.Now().UTC().Format(time.RFC3339),
			}
			def.Spec.RunDuration = "5s"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-declare-"
		},
		nil,
		"kubectl arcane downtime stop arcane-stream-mock-v2 downtime-window-1 --namespace integration-tests",
	)
}

func Test_DowntimeList(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Labels = map[string]string{
				interfaces.DowntimeLabelKey: "maintenance-window-1",
			}
			def.Annotations = map[string]string{
				interfaces.DowntimeBeginAnnotationKey: time.Now().UTC().Format(time.RFC3339),
			}
			def.Spec.RunDuration = "5s"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-list-"
		},
		nil,
		"kubectl arcane downtime list",
	)
}

func Test_DowntimeDetails(t *testing.T) {
	runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Labels = map[string]string{
				interfaces.DowntimeLabelKey: "maintenance-window-1",
			}
			def.Annotations = map[string]string{
				interfaces.DowntimeBeginAnnotationKey: time.Now().UTC().Format(time.RFC3339),
			}
			def.Spec.RunDuration = "5s"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-downtime-details-"
		},
		nil,
		"kubectl arcane downtime details",
	)
}

func Test_BackfillOverride(t *testing.T) {
	name := runIntegrationTest(t,
		func(def *v2.TestStreamDefinitionV2) {
			def.Namespace = "integration-tests"
			def.Spec.RunDuration = "10m"
			def.Spec.ExecutionSettings.Suspended = true
			def.Spec.ShouldFail = false
			def.GenerateName = "integration-test-backfill-override-"
		},
		func(ctx context.Context, clientSet *mockversionedv1.Clientset, name string, namespace string) error {
			err := helpers.WaitForPhase(t, clientSet, name, namespace, stream.Suspended)
			if err != nil {
				return err
			}
			return helpers.UnsuspendTestStreamDefinition(ctx, clientSet, name, namespace)
		},
		"kubectl arcane stream backfill arcane-stream-mock-v2 %s --namespace integration-tests --override .spec.shouldFail=true",
	)
	var job *batchv1.Job
	err := wait.PollUntilContextCancel(t.Context(), 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		job, err = kubernetesClientSet.BatchV1().Jobs("integration-tests").Get(ctx, name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return job.Labels["arcane/backfilling"] == "true", nil
	})
	require.NoError(t, err)

	t.Logf("Job name %s", job.Name)
	for env := range job.Spec.Template.Spec.Containers[0].Env {
		t.Logf("Found env variable %s with value %s", job.Spec.Template.Spec.Containers[0].Env[env].Name, job.Spec.Template.Spec.Containers[0].Env[env].Value)
		if job.Spec.Template.Spec.Containers[0].Env[env].Name == "STREAMCONTEXT__OVERRIDE" {
			return
		}
	}
	require.Fail(t, "STREAMCONTEXT__OVERRIDE was not found in job")
}

var (
	clientSet           *mockversionedv1.Clientset
	kubeconfigCmd       string
	kubernetesClientSet *kubernetes.Clientset
)

func TestMain(m *testing.M) {
	flag.StringVar(&kubeconfigCmd, "kubeconfig-cmd", "/opt/homebrew/bin/kind get kubeconfig", "Command to execute that outputs kubeconfig YAML content")
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
	if err != nil {
		panic(fmt.Errorf("error building clientSet: %w", err))
	}

	kubernetesClientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(fmt.Errorf("error building kubernetesClientSet: %w", err))
	}

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

func runIntegrationTest(t *testing.T, setup func(def *v2.TestStreamDefinitionV2), before func(context.Context, *mockversionedv1.Clientset, string, string) error, commandTemplate string) string {
	name := helpers.NewTestStream(t, clientSet, setup)
	require.NotEmpty(t, name)

	var command string
	if strings.Contains(commandTemplate, "%s") {
		command = fmt.Sprintf(commandTemplate, name)
	} else {
		command = commandTemplate
	}
	fmt.Println(command)

	if before != nil {
		require.NoError(t, before(t.Context(), clientSet, name, "integration-tests"))
	}

	output, err := runCommand(t.Context(), command)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
	}
	t.Logf("Command output:\n%s", string(output))

	return name
}
