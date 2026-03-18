package filter

import (
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services/interfaces"
)

var _ interfaces.ObjectFilter = (*AllowAll)(nil)

type AllowAll struct{}

func NewAllowAll() *AllowAll {
	return &AllowAll{}
}

func (a AllowAll) Matches(_ streamapis.Definition) (bool, error) {
	return true, nil
}
