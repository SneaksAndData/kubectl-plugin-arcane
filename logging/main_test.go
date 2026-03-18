package logging

import (
	"flag"
	"testing"
)

var (
	kubeconfigCmd string
)

// TestMain initializes the Kubernetes client and creates a suspended StreamDefinition for testing
func TestMain(m *testing.M) {
	// Initialize flat to avoid failures in the pipeline
	flag.StringVar(&kubeconfigCmd,
		"kubeconfig-cmd",
		"/opt/homebrew/bin/kind get kubeconfig",
		"Command to execute that outputs kubeconfig YAML content")
	flag.Parse()

	m.Run()
}
