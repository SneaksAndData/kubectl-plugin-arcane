package services

import (
	"fmt"
	"github.com/SneaksAndData/arcane-operator/pkg/apis/streaming/v1"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
)

type PrintableObject interface {
	runtime.Object
	apiv1.Object
}

func Printer(operation string) printers.ResourcePrinter {
	_ = v1.AddToScheme(scheme.Scheme)
	return printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.NamePrinter{ShortOutput: false, Operation: operation})
}

func FormatName(object PrintableObject) string {
	groupKind := printers.GetObjectGroupKind(object)
	return fmt.Sprintf("%s.%s/%s/%s", strings.ToLower(groupKind.Kind), groupKind.Group, object.GetNamespace(), object.GetName())
}
