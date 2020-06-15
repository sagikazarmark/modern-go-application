package command

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
)

type listOptions struct {
	client todov1.TodoListServiceClient
}

// NewListCommand creates a new cobra.Command for listing todo items.
func NewListCommand(c Context) *cobra.Command {
	options := listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List todo items",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			options.client = c.GetTodoClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runList(options)
		},
	}
	cobra.OnInitialize()

	return cmd
}

func runList(options listOptions) error {
	req := &todov1.ListItemsRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := options.client.ListItems(ctx, req)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Completed"})

	for _, item := range resp.GetItems() {
		table.Append([]string{item.GetId(), item.GetTitle(), strconv.FormatBool(item.GetCompleted())})
	}
	table.Render()

	return nil
}
