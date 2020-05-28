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

func (r *mutationResolver) AddTodoItem(ctx context.Context, input graphql.NewTodoItem) (string, error) {
	req := AddItemRequest{
		NewItem: todo.NewItem{Title: input.Title},
	}

	resp, err := r.endpoints.AddItem(ctx, req)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return "", errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return "", f.Failed()
	}

	return resp.(AddItemResponse).Item.ID, nil
}

func (r *mutationResolver) MarkTodoAsComplete(ctx context.Context, input string) (bool, error) {
	req := MarkAsCompleteRequest{
		Id: input,
	}

	resp, err := r.endpoints.MarkAsComplete(ctx, req)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return false, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return false, f.Failed()
	}

	return true, nil
}

type queryResolver struct{ *resolver }

func (r *queryResolver) TodoItems(ctx context.Context) ([]*todo.Item, error) {
	resp, err := r.endpoints.ListItems(ctx, nil)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return nil, errors.New("internal server error")
	}

	todos := make([]*todo.Item, len(resp.(ListItemsResponse).Items))

	for i, todo := range resp.(ListItemsResponse).Items {
		todo := todo
		todos[i] = &todo
	}

	return todos, nil
}
