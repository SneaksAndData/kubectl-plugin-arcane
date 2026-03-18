package interfaces

import (
	v1 "github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	"github.com/SneaksAndData/arcane-operator/services/controllers/stream"
)

type QueueItem struct {
	Definition stream.Definition
	Class      *v1.StreamClass
}
