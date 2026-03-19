package filter

import (
	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
)

var _ interfaces.ObjectFilter = (*ByDowntimeKey)(nil)

type ByDowntimeKey struct {
	key string
}

func NewByDowntimeKey(key string) *ByDowntimeKey {
	return &ByDowntimeKey{
		key: key,
	}
}

func (f *ByDowntimeKey) Matches(definition stream.Definition) (bool, error) {
	return definition.ToUnstructured().GetLabels()[interfaces.DowntimeAnnotationKey] == f.key, nil
}
