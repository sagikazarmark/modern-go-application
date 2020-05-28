package tododriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-kit/kit/endpoint"

	"github.com/sagikazarmark/modern-go-application/.gen/api/graphql"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

// MakeGraphQLHandler mounts all of the service endpoints into a GraphQL handler.
func MakeGraphQLHandler(endpoints Endpoints, errorHandler todo.ErrorHandler) http.Handler {
	return handler.New(
		graphql.NewExecutableSchema(graphql.Config{
			Resolvers: &resolver{
				endpoints:    endpoints,
				errorHandler: errorHandler,
			},
		}),
	)
}

type resolver struct {
	endpoints    Endpoints
	errorHandler todo.ErrorHandler
}

func (r *resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}
func (r *resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *resolver }

func (r *mutationResolver) AddTodoItem(ctx context.Context, input graphql.NewTodoItem) (*todo.Item, error) {
	var order int
	if input.Order != nil {
		order = *input.Order
	}

	req := AddItemRequest{
		NewItem: todo.NewItem{
			Title: input.Title,
			Order: order,
		},
	}

	resp, err := r.endpoints.AddItem(ctx, req)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return nil, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return nil, f.Failed()
	}

	item := resp.(AddItemResponse).Item

	return &item, nil
}

func (r *mutationResolver) UpdateTodoItem(ctx context.Context, input graphql.TodoItemUpdate) (*todo.Item, error) {
	req := UpdateItemRequest{
		Id: input.ID,
		ItemUpdate: todo.ItemUpdate{
			Title:     input.Title,
			Completed: input.Completed,
			Order:     input.Order,
		},
	}

	resp, err := r.endpoints.AddItem(ctx, req)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return nil, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return nil, f.Failed()
	}

	item := resp.(AddItemResponse).Item

	return &item, nil
}

type queryResolver struct{ *resolver }

func (r *queryResolver) TodoItems(ctx context.Context) ([]*todo.Item, error) {
	resp, err := r.endpoints.ListItems(ctx, nil)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return nil, errors.New("internal server error")
	}

	items := make([]*todo.Item, len(resp.(ListItemsResponse).Items))

	for i, item := range resp.(ListItemsResponse).Items {
		item := item
		items[i] = &item
	}

	return items, nil
}
