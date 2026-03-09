package services

import (
	"fmt"
	"os"
)

func logError(object PrintableObject, operation string, cause error) {
	name := FormatName(object)
	_, err := fmt.Fprintf(os.Stderr, "%s Failed %s: %v\n", name, operation, cause)
	if err != nil {
		panic(err)
	}
}
