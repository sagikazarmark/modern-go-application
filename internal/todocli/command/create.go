package command

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

type createOptions struct {
	text   string
	client todov1beta1.TodoListClient
}

// NewCreateCommand creates a new cobra.Command for creating a todo.
func NewCreateCommand(c Context) *cobra.Command {
	options := createOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a TODO",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.text = args[0]
			options.client = c.GetTodoClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runCreate(options)
		},
	}

	return cmd
}

func runCreate(options createOptions) error {
	req := &todov1beta1.CreateTodoRequest{
		Text: options.text,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := options.client.CreateTodo(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("Todo %q with ID %s has been created.", options.text, resp.GetId())

	return nil
}
