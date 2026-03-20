package publisher

import (
	"context"

	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/logging"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ interfaces.QueuePublisher = (*StreamClassMembers)(nil)

type StreamClassMembers struct {
	clientProvider cmdinterfaces.ClientProvider
	streamClass    string
	namespace      string
	objectFilter   interfaces.ObjectFilter
	selector       *client.MatchingLabelsSelector
}

func NewStreamClassMembersPublisher(provider cmdinterfaces.ClientProvider, streamClass string, namespace string, objectFilter interfaces.ObjectFilter, selector *client.MatchingLabelsSelector) *StreamClassMembers {
	return &StreamClassMembers{
		clientProvider: provider,
		streamClass:    streamClass,
		namespace:      namespace,
		objectFilter:   objectFilter,
		selector:       selector,
	}
}

func (s StreamClassMembers) PublishStreamDefinitions(ctx context.Context, queue interfaces.Queue) error {
	clientSet, err := s.clientProvider.ProvideClientSet()
	if err != nil { // coverage-ignore
		return err
	}
	sc, err := clientSet.
		StreamingV1().
		StreamClasses(""). // StreamClasses are cluster-scoped, so we ignore the namespace parameter here.
		Get(ctx, s.streamClass, metav1.GetOptions{})
	if err != nil { // coverage-ignore
		return err
	}

	gvk := sc.TargetResourceGvk()

	streamList := &unstructured.UnstructuredList{}
	streamList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind + "List",
	})

	unstructuredClient, err := s.clientProvider.ProvideUnstructuredClient()
	if err != nil { // coverage-ignore
		return err
	}

	err = unstructuredClient.List(ctx, streamList, client.InNamespace(s.namespace), s.selector)
	if err != nil { // coverage-ignore
		return err
	}

	for _, item := range streamList.Items {
		streamDefinition, err := streamapis.FromUnstructured(&item)
		if err != nil {
			logging.LogError(&item, "parsing kubernetes object, skipping", err)
			continue // Skip items that can't be parsed as stream definitions
		}

		matches, err := s.objectFilter.Matches(streamDefinition)
		if err != nil {
			logging.LogError(&item, "applying object filter, skipping", err)
			continue // Skip items that cause errors when applying the filter
		}
		if !matches {
			continue
		}
		queue.Add(interfaces.QueueItem{Definition: streamDefinition, Class: sc})
	}

	return nil
}
