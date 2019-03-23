package command

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
)

type listOptions struct {
	client todov1beta1.TodoListClient
}

// NewListCommand creates a new cobra.Command for listing todos.
func NewListCommand(c Context) *cobra.Command {
	options := listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List TODOs",
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
	req := &todov1beta1.ListTodosRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := options.client.ListTodos(ctx, req)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Task", "Done"})

	for _, t := range resp.GetTodos() {
		table.Append([]string{t.GetId(), t.GetText(), strconv.FormatBool(t.GetDone())})
	}
	table.Render()

	return nil
}
