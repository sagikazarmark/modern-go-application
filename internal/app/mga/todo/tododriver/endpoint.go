package tododriver

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

type createTodoRequest struct {
	Text string `json:"text"`
}

type createTodoResponse struct {
	ID string `json:"id"`
}

// MakeCreateTodoEndpoint returns an endpoint for the matching method of the underlying service.
func MakeCreateTodoEndpoint(service todo.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTodoRequest)

		id, err := service.CreateTodo(ctx, req.Text)

		return createTodoResponse{
			ID: id,
		}, err
	}
}

type listTodosResponse struct {
	Todos []todo.Todo `json:"todos"`
}

// MakeListTodosEndpoint returns an endpoint for the matching method of the underlying service.
func MakeListTodosEndpoint(service todo.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		todos, err := service.ListTodos(ctx)

		resp := listTodosResponse{
			Todos: todos,
		}

		return resp, err
	}
}

type markAsDoneRequest struct {
	ID string
}

// MakeMarkAsDoneEndpoint returns an endpoint for the matching method of the underlying service.
func MakeMarkAsDoneEndpoint(service todo.Service) endpoint.Endpoint {
	return kitxendpoint.BusinessErrorMiddleware(func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(markAsDoneRequest)

		return nil, service.MarkAsDone(ctx, req.ID)
	})
}
