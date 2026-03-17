package interfaces

import (
	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
)

// ObjectFilter defines an interface for filtering stream definitions based on specific criteria.
type ObjectFilter interface {
	// Matches determines whether the given stream definition matches the filter criteria.
	Matches(definition stream.Definition) (bool, error)
}
