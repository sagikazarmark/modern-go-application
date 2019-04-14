package tododriver

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-kit/kit/endpoint"

	"github.com/sagikazarmark/modern-go-application/.gen/api/graphql"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// MakeGraphQLHandler mounts all of the service endpoints into a GraphQL handler.
func MakeGraphQLHandler(endpoints Endpoints, errorHandler todo.ErrorHandler) http.Handler {
	return handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &resolver{
			endpoints:    endpoints,
			errorHandler: errorHandler,
		},
	}))
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
	req := createTodoRequest{
		Text: input.Text,
	}

	resp, err := r.endpoints.Create(ctx, req)
	if err != nil {
		r.errorHandler.Handle(err)

		return "", errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return "", f.Failed()
	}

	return resp.(createTodoResponse).ID, nil
}

func (r *mutationResolver) MarkTodoAsDone(ctx context.Context, input string) (bool, error) {
	req := markAsDoneRequest{
		ID: input,
	}

	resp, err := r.endpoints.MarkAsDone(ctx, req)
	if err != nil {
		r.errorHandler.Handle(err)

		return false, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok {
		return false, f.Failed()
	}

	return true, nil
}

type queryResolver struct{ *resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]graphql.Todo, error) {
	resp, err := r.endpoints.List(ctx, nil)
	if err != nil {
		r.errorHandler.Handle(err)

		return nil, errors.New("internal server error")
	}

	todoResp := resp.(listTodosResponse)

	todos := make([]graphql.Todo, len(todoResp.Todos))

	for i, t := range todoResp.Todos {
		todos[i] = graphql.Todo{
			ID:   t.ID,
			Text: t.Text,
			Done: t.Done,
		}
	}

	return todos, nil
}
