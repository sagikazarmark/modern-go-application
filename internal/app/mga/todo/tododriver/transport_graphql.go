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

func (r *mutationResolver) CreateTodo(ctx context.Context, input graphql.NewTodo) (string, error) {
	req := CreateTodoRequest{
		Text: input.Text,
	}

	resp, err := r.endpoints.CreateTodo(ctx, req)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return "", errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return "", f.Failed()
	}

	return resp.(CreateTodoResponse).Id, nil
}

func (r *mutationResolver) MarkTodoAsDone(ctx context.Context, input string) (bool, error) {
	req := MarkAsDoneRequest{
		Id: input,
	}

	resp, err := r.endpoints.MarkAsDone(ctx, req)
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

func (r *queryResolver) Todos(ctx context.Context) ([]*todo.Todo, error) {
	resp, err := r.endpoints.ListTodos(ctx, nil)
	if err != nil {
		r.errorHandler.HandleContext(ctx, err)

		return nil, errors.New("internal server error")
	}

	todos := make([]*todo.Todo, len(resp.(ListTodosResponse).Todos))

	for i, todo := range resp.(ListTodosResponse).Todos {
		todo := todo
		todos[i] = &todo
	}

	return todos, nil
}
