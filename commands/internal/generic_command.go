package internal

import "github.com/spf13/cobra"

type GenericCommand interface {
	GetCommand() *cobra.Command
}

type genericCommand struct {
	*cobra.Command
}

func (s *genericCommand) GetCommand() *cobra.Command {
	return s.Command
}

func NewGenericCommand(command *cobra.Command) GenericCommand {
	return &genericCommand{Command: command}
}
