package command

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

type markAsCompleteOptions struct {
	todoID string
	client todov1beta1.TodoListClient
}

// NewMarkAsCompleteCommand creates a new cobra.Command for marking a todo as complete.
func NewMarkAsCompleteCommand(c Context) *cobra.Command {
	options := markAsCompleteOptions{}

	cmd := &cobra.Command{
		Use:     "complete",
		Aliases: []string{"co"},
		Short:   "Mark a TODO as complete",
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
	req := &todov1beta1.MarkAsCompleteRequest{
		Id: options.todoID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := options.client.MarkAsComplete(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("Todo with ID %s has been marked as complete.", options.todoID)

	return nil
}
