package command

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

type markAsDoneOptions struct {
	todoID string
	client todov1beta1.TodoListClient
}

// NewMarkAsDoneCommand creates a new cobra.Command for marking a todo as done.
func NewMarkAsDoneCommand(c Context) *cobra.Command {
	options := markAsDoneOptions{}

	cmd := &cobra.Command{
		Use:     "done",
		Aliases: []string{"d"},
		Short:   "Mark a TODO as done",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.todoID = args[0]
			options.client = c.GetTodoClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runMarkAsDone(options)
		},
	}

	return cmd
}

func runMarkAsDone(options markAsDoneOptions) error {
	req := &todov1beta1.MarkAsDoneRequest{
		Id: options.todoID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := options.client.MarkAsDone(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("Todo with ID %s has been marked as done.", options.todoID)

	return nil
}
