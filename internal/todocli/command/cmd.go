package command

import (
	"github.com/spf13/cobra"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

// Context represents the application context.
type Context interface {
	GetTodoClient() todov1beta1.TodoListClient
}

// AddCommands adds all the commands from cli/command to the root command.
func AddCommands(cmd *cobra.Command, c Context) {
	cmd.AddCommand(
		NewCreateCommand(c),
		NewListCommand(c),
		NewMarkAsDoneCommand(c),
	)
}
