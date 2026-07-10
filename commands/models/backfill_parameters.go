package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// BackfillParameters represents the parameters required to perform a backfill operation for a stream.
type BackfillParameters struct {
	StreamClass string   // The class of the stream to backfill.
	StreamId    string   // The unique identifier of the stream to backfill.
	Wait        bool     // Whether to wait for the backfill operation to complete before returning.
	Namespace   string   // The namespace in which the stream is located. If empty, the default namespace will be used.
	overrides   []string // List of overrides to apply to the backfill operation, in the format "key=value".
}

// NewBackfillParameters creates a new instance of BackfillParameters based on the provided command and arguments.
func NewBackfillParameters(cmd *cobra.Command, args []string, configFlags *genericclioptions.ConfigFlags, overrides []string) (*BackfillParameters, error) { // coverage-ignore (tested in integration tests)

	wait, err := cmd.Flags().GetBool("wait")
	if err != nil {
		return nil, err
	}

	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	bfr := &BackfillParameters{
		StreamClass: args[0],
		StreamId:    args[1],
		Wait:        wait,
		Namespace:   namespace,
		overrides:   overrides,
	}

	return bfr, nil
}

func (p BackfillParameters) ToBackfillRequest() *v1.BackfillRequest {
	return &v1.BackfillRequest{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-manual-", p.StreamId),
		},
		Spec: v1.BackfillRequestSpec{
			StreamClass: p.StreamClass,
			StreamId:    p.StreamId,
			Payload:     generatePayload(p.overrides),
		},
	}
}

func generatePayload(overrides []string) *runtime.RawExtension {
	nestedSpecMap := make(map[string]interface{})
	for _, kv := range overrides {
		parts := strings.SplitN(kv, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}

		if strings.HasPrefix(key, ".spec.") {
			cleanKey := strings.TrimPrefix(key, ".spec.")
			setNestedValue(nestedSpecMap, cleanKey, value)
		}
	}

	jsonBytes, err := json.Marshal(nestedSpecMap)
	if err != nil {
		fmt.Printf("Error building spec JSON: %v\n", err)
		os.Exit(1)
	}

	// 3. Inject raw bytes directly into the runtime.RawExtension
	return &runtime.RawExtension{Raw: jsonBytes}
}

func setNestedValue(m map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := m

	for i, part := range parts {
		if part == "" {
			continue
		}
		// If it's the final key in the path, assign the actual value
		if i == len(parts)-1 {
			current[part] = value
			return
		}

		// Otherwise, find or create the next nested map level
		next, exists := current[part]
		if !exists {
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		} else {
			if nestedMap, ok := next.(map[string]interface{}); ok {
				current = nestedMap
			} else {
				// Overwrite conflicting non-map values if structural definitions overlap
				newMap := make(map[string]interface{})
				current[part] = newMap
				current = newMap
			}
		}
	}
}
