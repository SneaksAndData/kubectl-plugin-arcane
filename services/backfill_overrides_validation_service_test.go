package services

import (
	"fmt"
	"testing"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/SneaksAndData/arcane-stream-mock/pkg/apis/streaming/v2"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/sneaksAndData/kubectl-plugin-arcane/tests/helpers"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_BackfillOverridesValidation_Backfill_WithOverridableParameters(t *testing.T) {
	name := helpers.NewTestStream(t, clientSet, func(def *v2.TestStreamDefinitionV2) {
		def.Spec.RunDuration = "5s"
		def.Spec.ExecutionSettings.Suspended = true
	})
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	streamClass, err := clientSet.StreamingV1().StreamClasses("").Get(t.Context(), testStreamClass, metav1.GetOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, streamClass.Spec.OverridableFields)

	overrides := []string{fmt.Sprintf("%s=true", streamClass.Spec.OverridableFields[0])}

	backfillService := newBackfillOverridesValidationService(NewFakeClientProvider(clientSet, nil))
	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: testStreamClass,
		Wait:        false,
		Overrides:   &overrides,
	})
	require.NoError(t, err)

	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.NoError(t, err)
	require.False(t, bfr.Spec.Completed)
}

func Test_BackfillOverridesValidation_Backfill_NonOverridableParameter(t *testing.T) {
	name := helpers.NewTestStream(t, clientSet, func(def *v2.TestStreamDefinitionV2) {
		def.Spec.RunDuration = "5s"
		def.Spec.ExecutionSettings.Suspended = true
	})
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	overrides := []string{".spec.notOverridableField=true"}

	backfillService := newBackfillOverridesValidationService(NewFakeClientProvider(clientSet, nil))
	err := backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: testStreamClass,
		Wait:        false,
		Overrides:   &overrides,
	})
	require.EqualError(t, err, "backfillOverridesValidationService: parameter is not overridable: .spec.notOverridableField")

	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.Error(t, err)
	require.Nil(t, bfr)
}

func Test_BackfillOverridesValidation_Backfill_InvalidOverrideFormat(t *testing.T) {
	name := helpers.NewTestStream(t, clientSet, func(def *v2.TestStreamDefinitionV2) {
		def.Spec.RunDuration = "5s"
		def.Spec.ExecutionSettings.Suspended = true
	})
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	overrides := []string{"spec.shouldFail=true"}

	backfillService := newBackfillOverridesValidationService(NewFakeClientProvider(clientSet, nil))
	err := backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: testStreamClass,
		Wait:        false,
		Overrides:   &overrides,
	})
	require.EqualError(t, err, "backfillOverridesValidationService: error checking prefix: key must start with dot")

	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.Error(t, err)
	require.Nil(t, bfr)
}
