package main

import (
	"errors"
	"fmt"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands/interfaces"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(commands.NewStreamCommand),
		fx.Provide(commands.NewDowntimeStopCommand),
		fx.Provide(commands.NewDowntimeCommand),
		fx.Provide(commands.NewDowntimeDeclareCommand),
		fx.Provide(services.NewDowntimeService),
		fx.Provide(commands.NewRootCommand),
		fx.Provide(commands.NewStreamStop),
		fx.Provide(commands.NewStreamStart),
		fx.Provide(commands.NewStreamBackfill),
		fx.NopLogger,
		fx.Provide(fx.Annotate(services.NewStreamService, fx.As(new(interfaces.StreamService)))),
		fx.Invoke(
			func(rootCmd commands.RootCommand, shutDowner fx.Shutdowner) error {
				err := rootCmd.Execute()
				defer func() {
					shErr := shutDowner.Shutdown()
					if shErr != nil {
						err = errors.Join(shErr)
					}
				}()

				if err != nil {
					return fmt.Errorf("failed to execute root command: %w", err)
				}

				return nil
			},
		),
	)

	app.Run()
}
