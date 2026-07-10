package models

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestNewBackfillParameters(t *testing.T) {
	cmd := &cobra.Command{Use: "backfill"}
	cmd.Flags().Bool("wait", false, "wait for completion")
	require.NoError(t, cmd.Flags().Set("wait", "true"))

	configFlags := genericclioptions.NewConfigFlags(false)
	overrides := []string{
		".spec.backfillBehavior=merge",
		".spec.sink.mergeServiceClient.connectionUrl=http://somewhere",
	}

	parameters, err := NewBackfillParameters(cmd, []string{"stream-class", "stream-id"}, configFlags, overrides)
	require.NoError(t, err)
	require.Equal(t, "stream-class", parameters.StreamClass)
	require.Equal(t, "stream-id", parameters.StreamId)
	require.True(t, parameters.Wait)
	require.Equal(t, overrides, parameters.overrides)
}

func TestBackfillParameters_ToBackfillRequest_WithOverrides(t *testing.T) {
	parameters := BackfillParameters{
		StreamClass: "stream-class",
		StreamId:    "stream-id",
		Namespace:   "default",
		overrides: []string{
			".spec.backfillBehavior=merge",
			".spec.sink.mergeServiceClient.connectionUrl=http://somewhere",
			".spec.somethingElse=something_else",
		},
	}

	request := parameters.ToBackfillRequest()
	require.Equal(t, "stream-id-manual-", request.GenerateName)
	require.Equal(t, "stream-class", request.Spec.StreamClass)
	require.Equal(t, "stream-id", request.Spec.StreamId)
	require.NotNil(t, request.Spec.Payload)
	require.NotEmpty(t, request.Spec.Payload.Raw)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(request.Spec.Payload.Raw, &payload))
	require.Equal(t, "merge", payload["backfillBehavior"])

	sink, ok := payload["sink"].(map[string]interface{})
	require.True(t, ok)
	mergeServiceClient, ok := sink["mergeServiceClient"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "http://somewhere", mergeServiceClient["connectionUrl"])
	require.Equal(t, "something_else", payload["somethingElse"])
}

func TestGeneratePayload_IgnoresNonSpecAndHandlesConflicts(t *testing.T) {
	payload := generatePayload([]string{
		".spec.a=value",
		".spec.a.b=nested-value",
		".spec.empty",
		".status.ignored=ignored",
		"spec.alsoIgnored=ignored",
	})

	require.NotNil(t, payload)
	require.NotEmpty(t, payload.Raw)

	var actual map[string]interface{}
	require.NoError(t, json.Unmarshal(payload.Raw, &actual))

	a, ok := actual["a"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "nested-value", a["b"])
	require.Equal(t, "", actual["empty"])
	_, hasStatus := actual["status"]
	require.False(t, hasStatus)
	_, hasNonSpec := actual["spec"]
	require.False(t, hasNonSpec)
}
