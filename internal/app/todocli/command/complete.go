package command

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cobra"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
)

type markAsCompleteOptions struct {
	todoID string
	client todov1.TodoListServiceClient
}

// NewMarkAsCompleteCommand creates a new cobra.Command for marking a todo item as complete.
func NewMarkAsCompleteCommand(c Context) *cobra.Command {
	options := markAsCompleteOptions{}

	cmd := &cobra.Command{
		Use:     "complete",
		Aliases: []string{"c"},
		Short:   "Mark a todo item as complete",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.todoID = args[0]
			options.client = c.GetTodoClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runMarkAsComplete(options)
		},
	}

	return cmd
}

func runMarkAsComplete(options markAsCompleteOptions) error {
	req := &todov1.UpdateItemRequest{
		Id: options.todoID,
		Completed: &wrappers.BoolValue{
			Value: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := options.client.UpdateItem(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("Todo item with ID %s has been marked as complete.", options.todoID)

	return nil
}
