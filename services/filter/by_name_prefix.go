package filter

import (
	"strings"

	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
)

var _ interfaces.ObjectFilter = (*UnsuspendedByNamePrefix)(nil)

type UnsuspendedByNamePrefix struct {
	Prefix string
}

func NewUnsuspendedByNamePrefix(prefix string) *UnsuspendedByNamePrefix {
	return &UnsuspendedByNamePrefix{
		Prefix: prefix,
	}
}

func (f *UnsuspendedByNamePrefix) Matches(definition stream.Definition) (bool, error) {
	return strings.HasPrefix(definition.ToUnstructured().GetName(), f.Prefix) && !definition.Suspended(), nil
}
