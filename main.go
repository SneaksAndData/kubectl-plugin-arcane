package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneaksAndData/kubectl-plugin-arcane/commands"
	"github.com/sneaksAndData/kubectl-plugin-arcane/services"
	"go.uber.org/fx"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	app := fx.New(
		fx.Supply(configFlags),
		fx.Provide(commands.NewStreamCommand),
		fx.Provide(commands.NewDowntimeStopCommand),
		fx.Provide(commands.NewDowntimeCommand),
		fx.Provide(commands.NewDowntimeDeclareCommand),
		fx.Provide(services.NewDowntimeService),
		fx.Provide(commands.NewRootCommand),
		fx.Provide(commands.NewStreamStop),
		fx.Provide(commands.NewStreamStart),
		fx.Provide(commands.NewStreamBackfill),
		//fx.NopLogger,
		fx.Provide(services.NewStreamService),
		fx.Provide(services.NewClientProvider),
		fx.Invoke(
			func(rootCmd commands.RootCommand, shutDowner fx.Shutdowner, lifeCycle fx.Lifecycle) error {
				err := rootCmd.GetCommand().ExecuteContext(context.TODO())
				fmt.Println(err)
				defer func() {
					shErr := shutDowner.Shutdown(fx.ExitCode(0))
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
