package services

import (
	"testing"

	versionedv1 "github.com/SneaksAndData/arcane-operator/pkg/generated/clientset/versioned"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/models"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_Backfill_existing_stream_definition(t *testing.T) {
	name := createTestStreamDefinition(t, false, "5s", true)
	require.NotEmpty(t, name)

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	backfillService := NewValidatedBackfillService(NewFakeClientProvider(clientSet, c))
	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    name,
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.NoError(t, err)
	bfr, err := findBackfillRequestByName(t.Context(), "default", name)
	require.NoError(t, err)
	require.False(t, bfr.Spec.Completed)
}

func Test_Backfill_no_stream_definition(t *testing.T) {

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	backfillService := NewValidatedBackfillService(NewFakeClientProvider(clientSet, c))
	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    "invalid-stream-id",
		StreamClass: "arcane-stream-mock",
		Wait:        false,
	})
	require.EqualError(t, err, "error fetching stream definition: teststreamdefinitions.streaming.sneaksanddata.com \"invalid-stream-id\" not found")
	bfr, err := findBackfillRequestByName(t.Context(), "default", "invalid-stream-id")
	require.Error(t, err)
	require.Nil(t, bfr)
}

func Test_Backfill_no_stream_class(t *testing.T) {

	clientSet := versionedv1.NewForConfigOrDie(kubeConfig)
	c, err := client.New(kubeConfig, client.Options{})
	require.NoError(t, err)

	backfillService := NewValidatedBackfillService(NewFakeClientProvider(clientSet, c))
	err = backfillService.Backfill(t.Context(), &models.BackfillParameters{
		Namespace:   "default",
		StreamId:    "invalid-stream-id",
		StreamClass: "invalid-stream-class",
		Wait:        false,
	})
	require.EqualError(t, err, "validatedBackfill: error getting stream class: streamclasses.streaming.sneaksanddata.com \"invalid-stream-class\" not found")
	bfr, err := findBackfillRequestByName(t.Context(), "default", "invalid-stream-id")
	require.Error(t, err)
	require.Nil(t, bfr)
}
