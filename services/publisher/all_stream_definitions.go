package publisher

import (
	"context"

	cmdinterfaces "github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/filter"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ interfaces.QueuePublisher = (*AllStreamDefinitions)(nil)

type AllStreamDefinitions struct {
	provider cmdinterfaces.ClientProvider
	selector *pkgclient.MatchingLabelsSelector
}

func NewAllStreamDefinitionsPublisher(provider cmdinterfaces.ClientProvider, selector *pkgclient.MatchingLabelsSelector) *AllStreamDefinitions {
	return &AllStreamDefinitions{
		provider: provider,
		selector: selector,
	}
}

func (a AllStreamDefinitions) PublishStreamDefinitions(ctx context.Context, target interfaces.Queue) error {
	client, err := a.provider.ProvideClientSet()
	if err != nil {
		return err
	}

	streamClasses, err := client.StreamingV1().StreamClasses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, sc := range streamClasses.Items {
		queuePublisher := NewStreamClassMembersPublisher(a.provider, sc.Name, "", filter.NewAllowAll(), a.selector)
		err = queuePublisher.PublishStreamDefinitions(ctx, target)
		if err != nil {
			return err
		}
	}

	return nil
}
