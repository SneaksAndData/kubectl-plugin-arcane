package errors

import (
	"fmt"
	streamapis "github.com/SneaksAndData/arcane-operator/services/controllers/stream"
	"k8s.io/apimachinery/pkg/types"
)

type StatusNoOpError struct {
	Phase streamapis.Phase
	name  types.NamespacedName
}

// NewStatusNoOpError creates a new instance of StatusNoOpError with the provided phase and namespaced name.
func NewStatusNoOpError(phase streamapis.Phase, name types.NamespacedName) *StatusNoOpError {
	return &StatusNoOpError{
		Phase: phase,
		name:  name,
	}
}

// Error returns a string representation of the StatusNoOpError, indicating that the stream already has the desired phase.
func (e *StatusNoOpError) Error() string {
	return fmt.Sprintf("Stream already has desired phase %s: %s/%s", e.Phase, e.name.Namespace, e.name.Name)
}
