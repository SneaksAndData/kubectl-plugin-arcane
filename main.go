package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hello from kubectl-plugin-arcane!\n")
		fmt.Printf("Arguments: %v\n", os.Args[1:])
	} else {
		fmt.Println("kubectl-plugin-arcane")
	}
}
