package command

import (
	"github.com/spf13/cobra"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
)

// Context represents the application context.
type Context interface {
	GetTodoClient() todov1.TodoListServiceClient
}

// AddCommands adds all the commands from cli/command to the root command.
func AddCommands(cmd *cobra.Command, c Context) {
	cmd.AddCommand(
		NewAddCommand(c),
		NewListCommand(c),
		NewMarkAsCompleteCommand(c),
	)
}
