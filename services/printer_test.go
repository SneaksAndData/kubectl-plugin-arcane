package services

import (
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func Test_FormatName(t *testing.T) {
	// Arrange
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)
	object := v1.BackfillRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-name",
			Namespace: "example-namespace",
		},
	}
	gvks, _, err := scheme.ObjectKinds(&object)
	require.NoError(t, err)
	require.Len(t, gvks, 1)
	object.SetGroupVersionKind(gvks[0])

	/// Act
	str := FormatName(&object)

	// Assert
	require.Equal(t, "backfillrequest.streaming.sneaksanddata.com/example-namespace/example-name", str)
}
