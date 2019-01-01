package workflow

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App struct {
	logger Logger
}

func NewApp(logger Logger) App {
	return App{logger}
}

func (app App) Run(ctx context.Context, args []string, signal <-chan os.Signal) error {
	cmd := &cobra.Command{
		Use: "workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("require command")
		},
	}
	cmd.SetArgs(args[1:])

	file := cmd.PersistentFlags().StringP("file", "f", "workflow.yaml", "workflow definition file.")

	// parse --file flag before execute to parse and append commands.
	if err := cmd.PersistentFlags().Parse(args); err != nil && err != pflag.ErrHelp {
		cmd.Usage()
		return err
	}

	def, err := LoadDefinition(*file)
	if err != nil {
		cmd.Usage()
		return err
	}

	for _, c := range def.Commands {
		cmd.AddCommand(&cobra.Command{
			Use: c.Name,
			RunE: func(cmd *cobra.Command, args []string) error {
				command, ok := def.Commands[cmd.Use]
				if !ok {
					return fmt.Errorf("no such workflow: %s", cmd.Use)
				}

				supervisor := NewSupervisor(app.logger)
				dispatcher := NewDispatcher(app.logger)

				var initialRunners []*TaskRunner
				for _, task := range command.Workflow {
					tr := NewTaskRunner(task, def.Env, supervisor, dispatcher)
					if len(task.When) == 0 {
						initialRunners = append(initialRunners, tr)
					} else {
						dispatcher.Register(ctx, tr)
					}
				}

				for _, tr := range initialRunners {
					tr.Run(ctx)
				}

				return supervisor.Wait()
			},
		})
	}

	return cmd.Execute()
}
