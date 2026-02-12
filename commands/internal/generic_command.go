package internal

import "github.com/spf13/cobra"

// GenericCommand is a wrapper around the Cobra Command struct, which allows us to inject commands to DI container.
type GenericCommand interface {
	GetCommand() *cobra.Command
}

// genericCommand is a simple implementation of the GenericCommand interface, which wraps a Cobra Command.
type genericCommand struct {
	*cobra.Command
}

// GetCommand returns the Cobra Command wrapped by the genericCommand struct.
func (s *genericCommand) GetCommand() *cobra.Command { // coverage-ignore (trivial)
	return s.Command
}

// NewGenericCommand creates a new instance of the GenericCommand interface, wrapping the provided Cobra Command.
func NewGenericCommand(command *cobra.Command) GenericCommand { // coverage-ignore (trivial)
	return &genericCommand{Command: command}
}
