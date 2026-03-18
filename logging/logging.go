package logging

import (
	"fmt"
	"os"
)

func LogError(object PrintableObject, operation string, cause error) { // coverage-ignore
	name := FormatName(object)
	_, err := fmt.Fprintf(os.Stderr, "%s Failed %s: %v\n", name, operation, cause)
	if err != nil {
		panic(err)
	}
}
